// Package ports define las interfaces (puertos) que establecen contratos entre
// las diferentes capas de la aplicación siguiendo la arquitectura hexagonal.
//
// mapper.go define Mapperport, la interfaz que establece el contrato para
// cualquier componente mapper. Los mappers son responsables de transformar
// datos crudos en objetos DTO estructurados. Esta abstracción permite
// diferentes implementaciones de mapeo según la fuente de datos.
package ports

import "github.com/dst3v3n/api-anime/internal/domain/dto"

// Mapperport define el contrato que debe cumplir cualquier implementación de mapper.
// Proporciona métodos para transformar datos crudos en estructuras DTO bien definidas.
type Mapperport interface {
	// ToSearchanime transforma datos básicos de búsqueda en un DTO AnimeResponse.
	ToSearchanime(id string, title string, sipnopsis string, tipo string, puctuation float64, image string) dto.AnimeResponse

	// ToAnimeinfo transforma datos completos de anime en un DTO AnimeInfoResponse.
	ToAnimeinfo(id string, title string, sipnopsis string, tipo string, puctuation float64, image string, animerelated []dto.AnimeRelated, generos []string, estado string, episodes []int, nextepisode string) dto.AnimeInfoResponse
}
