package testport

import (
	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"
)

type ScraperTestContract interface {
	SearchAnime(anime string, page string) ([]dto.AnimeResponse, error)
}
