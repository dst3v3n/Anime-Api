// Package animeflv implementa un cliente scraper para el sitio web AnimeFlv.
// Este archivo (client.go) contiene la estructura principal del cliente y la implementación
// de todos los métodos definidos en el port ScraperPort. Se encarga de realizar las
// peticiones HTTP a AnimeFlv y delegar el parsing del HTML al componente Parser.
package animeflv

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/ports"
	"golang.org/x/time/rate"
)

// Client es la estructura principal del scraper de AnimeFlv.
// Contiene la configuración de URLs del sitio y una instancia del parser HTML.
type Client struct {
	config  Config
	parser  *Parser
	limiter *rate.Limiter
	client  *http.Client
}

// NewClient crea una nueva instancia del cliente scraper de AnimeFlv.
// Inicializa la configuración con las URLs del sitio y crea el parser HTML.
// Retorna una interfaz ScraperPort para permitir la inyección de dependencias.
func NewClient() ports.ScraperPort {
	return &Client{
		config: Config{
			BaseURL:       "https://www3.animeflv.net",
			SearchURL:     "https://www3.animeflv.net/browse",
			AnimeInfoURL:  "https://www3.animeflv.net/anime",
			VerEpisodeURL: "https://www3.animeflv.net/ver",
		},
		parser:  NewParser(),
		limiter: rate.NewLimiter(rate.Limit(3), 5),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest es el método centralizado para realizar todas las peticiones HTTP.
// Aplica rate limiting automáticamente antes de cada petición y maneja timeouts.
// Este método garantiza que todas las peticiones al sitio respeten los límites establecidos.
func (c *Client) doRequest(ctx context.Context, url string) (*http.Response, error) {
	// Espera hasta que el rate limiter permita la petición
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter cancelado: %w", err)
	}

	// Crea la petición HTTP con el contexto
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición HTTP: %w", err)
	}

	// Realiza la petición
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en la petición HTTP: %w", err)
	}

	// Valida el código de estado
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("código de estado HTTP inesperado: %d", resp.StatusCode)
	}

	return resp, nil
}

// SearchAnime busca animes por nombre con soporte de paginación.
// Realiza una petición HTTP GET al endpoint de búsqueda de AnimeFlv
// y delega el parsing del HTML al componente Parser.
func (c *Client) SearchAnime(ctx context.Context, anime string, page string) (dto.AnimeResponse, error) {
	if page == "" {
		page = "1"
	}
	params := map[string]string{
		"page": page,
		"q":    anime,
	}
	searchURL := buildURL(c.config.SearchURL, params)

	resp, err := c.doRequest(ctx, searchURL)
	if err != nil {
		return dto.AnimeResponse{}, err
	}

	defer resp.Body.Close()

	return c.parser.ParseAnimeWithPagination(resp.Body)
}

// Search obtiene la lista de todos los animes disponibles sin filtros de búsqueda.
// Realiza una petición HTTP GET a la página de búsqueda sin parámetros de consulta
// y retorna todos los animes con información de paginación.
func (c *Client) Search(ctx context.Context) (dto.AnimeResponse, error) {
	resp, err := c.doRequest(ctx, c.config.SearchURL)
	if err != nil {
		return dto.AnimeResponse{}, err
	}

	defer resp.Body.Close()

	return c.parser.ParseAnimeWithPagination(resp.Body)
}

// AnimeInfo obtiene información detallada de un anime específico por su ID.
// Incluye sinopsis completa, géneros, estado de emisión, episodios disponibles,
// animes relacionados y fecha del próximo episodio si aplica.
func (c *Client) AnimeInfo(ctx context.Context, idAnime string) (dto.AnimeInfoResponse, error) {
	searchURL := c.config.AnimeInfoURL + "/" + idAnime
	resp, err := c.doRequest(ctx, searchURL)
	if err != nil {
		return dto.AnimeInfoResponse{}, err
	}

	defer resp.Body.Close()

	return c.parser.ParseAnimeInfo(resp.Body, idAnime)
}

// Links obtiene los enlaces de reproducción/descarga de un episodio específico.
// Retorna información de múltiples servidores de video con sus URLs y códigos de embed.
func (c *Client) Links(ctx context.Context, idAnime string, episode uint) (dto.LinkResponse, error) {
	searchURL := fmt.Sprintf("%s/%s-%d", c.config.VerEpisodeURL, idAnime, episode)
	resp, err := c.doRequest(ctx, searchURL)
	if err != nil {
		return dto.LinkResponse{}, err
	}

	defer resp.Body.Close()

	return c.parser.ParseLinks(resp.Body, idAnime, episode)
}

// RecentAnime obtiene la lista de animes recientemente agregados al sitio.
// Scrapeando la página principal de AnimeFlv para extraer los últimos contenidos publicados.
func (c *Client) RecentAnime(ctx context.Context) ([]dto.AnimeStruct, error) {
	resp, err := c.doRequest(ctx, c.config.BaseURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return c.parser.ParseAnime(resp.Body)
}

// RecentEpisode obtiene la lista de episodios recientemente publicados.
// Incluye información del anime, número de episodio, capítulo e imagen de portada.
func (c *Client) RecentEpisode(ctx context.Context) ([]dto.EpisodeListResponse, error) {
	resp, err := c.doRequest(ctx, c.config.BaseURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return c.parser.ParseRecentEpisode(resp.Body)
}
