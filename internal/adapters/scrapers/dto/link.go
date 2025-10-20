// Package dto - link.go
// Este archivo define las estructuras para manejar enlaces de reproducción de episodios.
// LinkResponse contiene información de un episodio específico junto con todos sus
// enlaces de reproducción disponibles. LinkSource representa cada servidor de video
// individual con su URL y código de reproducción.
package dto

type LinkResponse struct {
	ID      string
	Title   string
	Episode int
	Link    []LinkSource
}

type LinkSource struct {
	Server string
	URL    string
	Code   string
}
