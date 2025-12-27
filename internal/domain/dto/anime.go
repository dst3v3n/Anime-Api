// Package dto contiene los Data Transfer Objects (DTOs) utilizados para transferir
// información de anime entre las diferentes capas de la aplicación.
//
// anime.go define la estructura básica de respuesta de anime (AnimeResponse) que incluye
// información fundamental como ID, título, sinopsis, tipo, puntuación e imagen.
// También define los tipos de categoría de anime disponibles (Anime, OVA, Película, Especial).
package dto

// CategoryAnime representa el tipo de categoría de contenido de anime.
type CategoryAnime string

const (
	Anime    CategoryAnime = "Anime"    // Serie de anime regular
	Ova      CategoryAnime = "Ova"      // Original Video Animation
	Pelicula CategoryAnime = "Pelicula" // Película de anime
	Especial CategoryAnime = "Especial" // Especial de anime
)

// AnimeStruct contiene la información básica de un anime.
type AnimeStruct struct {
	ID          string        // Identificador único del anime (ej: "one-piece-tv")
	Title       string        // Título del anime
	Sinopsis    string        // Sinopsis o descripción del anime
	Type        CategoryAnime // Tipo/Categoría del anime
	Punctuation float64       // Calificación/puntuación del anime
	Image       string        // URL de la imagen/carátula del anime
}

// AnimeResponse es la estructura de respuesta para búsquedas de animes.
// Contiene una lista de animes y información de paginación.
type AnimeResponse struct {
	Animes     []AnimeStruct // Lista de animes encontrados
	TotalPages uint          // Número total de páginas disponibles para paginación
}
