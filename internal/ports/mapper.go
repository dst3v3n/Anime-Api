package ports

import "github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"

type Mapperport interface {
	ToSearchanime(id string, title string, sipnopsis string, tipo string, puctuation float64, image string) dto.SearchAnimeResponse
	ToAnimeinfo(id string, title string, sipnopsis string, tipo string, puctuation float64, image string, animerelated []dto.AnimeRelated, generos []string, estado string, episodes []int, nextepisode string) dto.ResponseAnimeInfo
}
