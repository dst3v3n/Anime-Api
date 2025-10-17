package ports

import "github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"

// ScraperPort define el contrato que debe cumplir cualquier scraper de anime
type ScraperPort interface {
	SearchAnime(anime string, page string) ([]dto.SearchAnimeResponse, error)
	AnimeInfo(idAnime string) (dto.ResponseAnimeInfo, error)
	GetLinks(idAnime string, episode int) (dto.GetLinkResponse, error)
}
