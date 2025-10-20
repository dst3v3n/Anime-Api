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

type StatusAnime string

const (
	Emision    StatusAnime = "En Emision"
	Finalizado StatusAnime = "Finalizado"
)

type AnimeInfoResponse struct {
	AnimeResponse
	AnimeRelated []AnimeRelated
	Genres       []string
	Status       StatusAnime
	NextEpisode  string
	Episodes     []int
}

type AnimeRelated struct {
	ID       string
	Title    string
	Category string
}
