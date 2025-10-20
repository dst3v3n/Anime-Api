// Package dto - episodeList.go
// Este archivo define la estructura EpisodeListResponse utilizada para representar
// episodios en listados (como episodios recientes). Contiene información resumida
// de cada episodio: ID del anime, título, capítulo, número de episodio e imagen.
package dto

type EpisodeListResponse struct {
	ID      string
	Title   string
	Chapter string
	Episode int
	Image   string
}
