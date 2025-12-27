// Package dto - episodeList.go
// Este archivo define la estructura EpisodeListResponse utilizada para representar
// episodios en listados (como episodios recientes). Contiene información resumida
// de cada episodio: ID del anime, título, capítulo, número de episodio e imagen.
package dto

// EpisodeListResponse contiene la información resumida de un episodio en un listado.
// Se utiliza para mostrar episodios recientes u otros listados de episodios sin detalles completos.
type EpisodeListResponse struct {
	ID      string // Identificador único del anime
	Title   string // Título del anime
	Chapter string // Designación del capítulo (ej: "Cap. 1050")
	Episode int    // Número del episodio
	Image   string // URL de la imagen/carátula del episodio
}
