// Package ports - mapper.go
// Este archivo define Mapperport, la interfaz que establece el contrato para
// cualquier componente mapper. Los mappers son responsables de transformar
// datos crudos en objetos DTO estructurados. Esta abstracción permite
// diferentes implementaciones de mapeo según la fuente de datos.
package ports

import "github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"

type Mapperport interface {
	ToSearchanime(id string, title string, sipnopsis string, tipo string, puctuation float64, image string) dto.AnimeResponse
	ToAnimeinfo(id string, title string, sipnopsis string, tipo string, puctuation float64, image string, animerelated []dto.AnimeRelated, generos []string, estado string, episodes []int, nextepisode string) dto.AnimeInfoResponse
}
