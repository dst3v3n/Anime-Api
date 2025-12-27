// Package dto - anime_info.go
// Este archivo extiende la información básica de anime con detalles adicionales.
// Define AnimeInfoResponse que incluye información completa sobre un anime:
// - Datos básicos heredados de AnimeResponse
// - Animes relacionados (secuelas, precuelas, spin-offs)
// - Géneros del anime
// - Estado de emisión (En Emisión o Finalizado)
// - Información del próximo episodio
// - Lista completa de episodios disponibles
package dto

// StatusAnime representa el estado de emisión del anime.
type StatusAnime string

const (
	Emision    StatusAnime = "En Emision" // Anime actualmente en emisión
	Finalizado StatusAnime = "Finalizado" // Anime finalizado
)

// AnimeInfoResponse contiene información completa y detallada de un anime específico.
type AnimeInfoResponse struct {
	AnimeStruct                 // Información básica del anime
	AnimeRelated []AnimeRelated // Animes relacionados (secuelas, precuelas, spin-offs, etc.)
	Genres       []string       // Géneros del anime (Acción, Aventura, Romance, etc.)
	Status       StatusAnime    // Estado actual de emisión del anime
	NextEpisode  string         // Fecha del próximo episodio a emitirse
	Episodes     []int          // Lista de números de episodios disponibles
}

// AnimeRelated contiene información básica de animes relacionados.
type AnimeRelated struct {
	ID       string // Identificador único del anime relacionado
	Title    string // Título del anime relacionado
	Category string // Tipo de relación (Secuela, Precuela, Spin-off, etc.)
}
