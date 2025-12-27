// Package animeflv implementa un cliente scraper para el sitio web AnimeFlv.
// Este archivo (client.go) contiene la estructura principal del cliente y la implementación
// de todos los métodos definidos en el port ScraperPort. Se encarga de realizar las
// peticiones HTTP a AnimeFlv y delegar el parsing del HTML al componente Parser.
package animeflv

import (
	"fmt"
	"net/http"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/ports"
)

// Client es la estructura principal del scraper de AnimeFlv.
// Contiene la configuración de URLs del sitio y una instancia del parser HTML.
type Client struct {
	config Config
	parser *Parser
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
		parser: NewParser(),
	}
}

// SearchAnime busca animes por nombre con soporte de paginación.
// Realiza una petición HTTP GET al endpoint de búsqueda de AnimeFlv
// y delega el parsing del HTML al componente Parser.
func (c *Client) SearchAnime(anime string, page string) (dto.AnimeResponse, error) {
	if page == "" {
		page = "1"
	}
	params := map[string]string{
		"page": page,
		"q":    anime,
	}
	searchURL := buildURL(c.config.SearchURL, params)
	resp, err := http.Get(searchURL)
	if err != nil {
		return dto.AnimeResponse{}, fmt.Errorf("error en la petición HTTP de búsqueda: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.AnimeResponse{}, fmt.Errorf("código de estado HTTP inesperado: %d", resp.StatusCode)
	}

	return c.parser.ParseAnimeWithPagination(resp.Body)
}

// Search obtiene la lista de todos los animes disponibles sin filtros de búsqueda.
// Realiza una petición HTTP GET a la página de búsqueda sin parámetros de consulta
// y retorna todos los animes con información de paginación.
func (c *Client) Search() (dto.AnimeResponse, error) {
	resp, err := http.Get(c.config.SearchURL)
	if err != nil {
		return dto.AnimeResponse{}, fmt.Errorf("error en la petición HTTP de búsqueda: %w", err)
	}

	defer resp.Body.Close()

	return c.parser.ParseAnimeWithPagination(resp.Body)
}

// AnimeInfo obtiene información detallada de un anime específico por su ID.
// Incluye sinopsis completa, géneros, estado de emisión, episodios disponibles,
// animes relacionados y fecha del próximo episodio si aplica.
func (c *Client) AnimeInfo(idAnime string) (dto.AnimeInfoResponse, error) {
	searchURL := c.config.AnimeInfoURL + "/" + idAnime
	resp, err := http.Get(searchURL)
	if err != nil {
		return dto.AnimeInfoResponse{}, fmt.Errorf("error en la petición HTTP de información del anime: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return dto.AnimeInfoResponse{}, fmt.Errorf("código de estado HTTP inesperado: %d", resp.StatusCode)
	}

	return c.parser.ParseAnimeInfo(resp.Body, idAnime)
}

// Links obtiene los enlaces de reproducción/descarga de un episodio específico.
// Retorna información de múltiples servidores de video con sus URLs y códigos de embed.
func (c *Client) Links(idAnime string, episode uint) (dto.LinkResponse, error) {
	searchURL := fmt.Sprintf("%s/%s-%d", c.config.VerEpisodeURL, idAnime, episode)
	resp, err := http.Get(searchURL)
	if err != nil {
		return dto.LinkResponse{}, fmt.Errorf("error en la petición HTTP de enlaces del episodio: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.LinkResponse{}, fmt.Errorf("código de estado HTTP inesperado: %d", resp.StatusCode)
	}

	return c.parser.ParseLinks(resp.Body, idAnime, episode)
}

// RecentAnime obtiene la lista de animes recientemente agregados al sitio.
// Scrapeando la página principal de AnimeFlv para extraer los últimos contenidos publicados.
func (c *Client) RecentAnime() ([]dto.AnimeStruct, error) {
	searchURL := c.config.BaseURL
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("error en la petición HTTP de animes recientes: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("código de estado HTTP inesperado: %d", resp.StatusCode)
	}

	return c.parser.ParseAnime(resp.Body)
}

// RecentEpisode obtiene la lista de episodios recientemente publicados.
// Incluye información del anime, número de episodio, capítulo e imagen de portada.
func (c *Client) RecentEpisode() ([]dto.EpisodeListResponse, error) {
	searchURL := c.config.BaseURL
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("error en la petición HTTP de episodios recientes: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("código de estado HTTP inesperado: %d", resp.StatusCode)
	}

	return c.parser.ParseRecentEpisode(resp.Body)
}
