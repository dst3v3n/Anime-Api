// Package ports define las interfaces (puertos) que establecen contratos entre
// las diferentes capas de la aplicación siguiendo la arquitectura hexagonal.
//
// scraper.go define ScraperPort, la interfaz principal que debe implementar
// cualquier scraper de sitios de anime. Esto permite cambiar la fuente de datos
// (por ejemplo, de AnimeFlv a otro sitio) sin afectar la lógica de negocio.
// Define operaciones como búsqueda, información detallada, obtención de enlaces
// de reproducción, y listado de contenido reciente.
package ports

import "github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"

// ScraperPort define el contrato que debe cumplir cualquier scraper de anime
type ScraperPort interface {
	SearchAnime(anime string, page string) ([]dto.AnimeResponse, error)
	AnimeInfo(idAnime string) (dto.AnimeInfoResponse, error)
	Links(idAnime string, episode int) (dto.LinkResponse, error)
	RecentAnime() ([]dto.AnimeResponse, error)
	RecentEpisode() ([]dto.EpisodeListResponse, error)
}
