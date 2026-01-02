// Package ports define las interfaces (puertos) que establecen contratos entre
// las diferentes capas de la aplicación siguiendo la arquitectura hexagonal.
//
// scraper.go define ScraperPort, la interfaz principal que debe implementar
// cualquier scraper de sitios de anime. Esto permite cambiar la fuente de datos
// (por ejemplo, de AnimeFlv a otro sitio) sin afectar la lógica de negocio.
// Define operaciones como búsqueda, información detallada, obtención de enlaces
// de reproducción, y listado de contenido reciente.
package ports

import (
	"context"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
)

// ScraperPort define el contrato que debe cumplir cualquier scraper de anime
type ScraperPort interface {
	SearchAnime(ctx context.Context, anime string, page string) (dto.AnimeResponse, error)
	Search(ctx context.Context, page string) (dto.AnimeResponse, error)
	AnimeInfo(ctx context.Context, idAnime string) (dto.AnimeInfoResponse, error)
	Links(ctx context.Context, idAnime string, episode uint) (dto.LinkResponse, error)
	RecentAnime(ctx context.Context) ([]dto.AnimeStruct, error)
	RecentEpisode(ctx context.Context) ([]dto.EpisodeListResponse, error)
}
