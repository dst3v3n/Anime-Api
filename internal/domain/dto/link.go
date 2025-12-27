// Package dto - link.go
// Este archivo define las estructuras para manejar enlaces de reproducción de episodios.
// LinkResponse contiene información de un episodio específico junto con todos sus
// enlaces de reproducción disponibles. LinkSource representa cada servidor de video
// individual con su URL y código de reproducción.
package dto

// LinkResponse contiene información de un episodio con sus enlaces de reproducción disponibles.
type LinkResponse struct {
	ID      string       // Identificador único del anime
	Title   string       // Título del anime
	Episode uint         // Número del episodio
	Link    []LinkSource // Lista de enlaces de reproducción disponibles para este episodio
}

// LinkSource representa un servidor de video individual para reproducción.
type LinkSource struct {
	Server string // Nombre del servidor de video (ej: "Zippyshare", "Mega", "Google Drive")
	URL    string // URL del enlace de reproducción/descarga
	Code   string // Código de embed o identificador del video en el servidor
}
