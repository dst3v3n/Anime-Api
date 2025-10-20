// Package dto contiene los Data Transfer Objects (DTOs) utilizados para transferir
// información de anime entre las diferentes capas de la aplicación.
//
// anime.go define la estructura básica de respuesta de anime (AnimeResponse) que incluye
// información fundamental como ID, título, sinopsis, tipo, puntuación e imagen.
// También define los tipos de categoría de anime disponibles (Anime, OVA, Película, Especial).
package dto

type CategoryAnime string

const (
	Anime    CategoryAnime = "Anime"
	Ova      CategoryAnime = "Ova"
	Pelicula CategoryAnime = "Pelicula"
	Especial CategoryAnime = "Especial"
)

type AnimeResponse struct {
	ID          string
	Title       string
	Sinopsis    string
	Type        CategoryAnime
	Punctuation float64
	Image       string
}
