// Package apianime es la API pública de la librería.
// Proporciona una interfaz simplificada para acceder a todos los servicios de scraping de anime.
// Los usuarios externos importan este paquete para usar la funcionalidad de la librería.
package apianime

import (
	"context"

	"github.com/dst3v3n/api-anime/internal/domain/services/animeflv"
	"github.com/dst3v3n/api-anime/types"
)

// AnimeFlv es la fachada principal que expone públicamente todos los servicios de anime.
// Encapsula el servicio interno de dominio y proporciona métodos para búsqueda, información
// detallada, enlaces de reproducción y contenido reciente.
type AnimeFlv struct {
	service *animeflv.AnimeflvService
}

// NewAnimeFlv crea una nueva instancia del servicio público de AnimeFlv.
// Inicializa el servicio interno con todas sus dependencias (scraper, caché, etc.).
func NewAnimeFlv() *AnimeFlv {
	return &AnimeFlv{
		service: animeflv.NewAnimeflvService(),
	}
}

// SearchAnime busca animes por nombre con soporte de paginación.
// Delega la operación al servicio interno de búsqueda.
func (s *AnimeFlv) SearchAnime(ctx context.Context, anime string, page uint) (types.AnimeResponse, error) {
	return s.service.SearchAnime(ctx, anime, page)
}

// Search obtiene todos los animes disponibles sin filtros de búsqueda.
// Delega la operación al servicio interno de búsqueda.
func (s *AnimeFlv) Search(ctx context.Context) (types.AnimeResponse, error) {
	return s.service.Search(ctx)
}

// AnimeInfo obtiene información detallada de un anime por su ID.
// Delega la operación al servicio interno de detalles.
func (s *AnimeFlv) AnimeInfo(ctx context.Context, idAnime string) (types.AnimeInfoResponse, error) {
	return s.service.AnimeInfo(ctx, idAnime)
}

// Links obtiene los enlaces de reproducción para un episodio específico.
// Retorna información de múltiples servidores de video con URLs y códigos de embed.
func (s *AnimeFlv) Links(ctx context.Context, idAnime string, episode uint) (types.LinkResponse, error) {
	return s.service.Links(ctx, idAnime, episode)
}

// RecentAnime obtiene la lista de animes recientemente agregados al sitio.
func (s *AnimeFlv) RecentAnime(ctx context.Context) ([]types.AnimeStruct, error) {
	return s.service.RecentAnime(ctx)
}

// RecentEpisode obtiene la lista de episodios recientemente publicados.
func (s *AnimeFlv) RecentEpisode(ctx context.Context) ([]types.EpisodeListResponse, error) {
	return s.service.RecentEpisode(ctx)
}
