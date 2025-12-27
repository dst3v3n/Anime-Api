// Package animeflv - mapper.go
// Este archivo implementa el componente Mapper que se encarga de transformar
// datos crudos extraídos del HTML en estructuras DTO (Data Transfer Objects) bien definidas.
// Convierte los datos primitivos en objetos de dominio utilizables por la aplicación,
// proporcionando una capa de abstracción entre el scraping y la lógica de negocio.
package animeflv

import "github.com/dst3v3n/api-anime/internal/domain/dto"

// Maper es el componente encargado de transformar datos a DTOs.
type Maper struct{}

// NewMaper crea una nueva instancia del mapper.
func NewMaper() *Maper {
	return &Maper{}
}

// ToAnime transforma datos básicos de anime en un DTO AnimeResponse.
// Convierte tipos primitivos en una estructura bien definida.
func (m *Maper) ToAnime(ID string, Title string, Sinopsis string, Tipo string, Punctuation float64, Image string) dto.AnimeStruct {
	return dto.AnimeStruct{
		ID:          ID,
		Title:       Title,
		Sinopsis:    Sinopsis,
		Type:        dto.CategoryAnime(Tipo),
		Punctuation: Punctuation,
		Image:       Image,
	}
}

// ToAnimeInfo transforma datos completos de anime en un DTO AnimeInfoResponse.
// Combina información básica con datos adicionales como géneros, episodios y animes relacionados.
func (m *Maper) ToAnimeInfo(ID string, Title string, Sinopsis string, Tipo string, Punctuation float64, Image string, AnimeRelated []dto.AnimeRelated, Generos []string, Estado string, Episodes []int, NextEpisode string) dto.AnimeInfoResponse {
	return dto.AnimeInfoResponse{
		AnimeStruct: dto.AnimeStruct{
			ID:          ID,
			Title:       Title,
			Sinopsis:    Sinopsis,
			Type:        dto.CategoryAnime(Tipo),
			Punctuation: Punctuation,
			Image:       Image,
		},
		AnimeRelated: AnimeRelated,
		Genres:       Generos,
		Status:       dto.StatusAnime(Estado),
		NextEpisode:  NextEpisode,
		Episodes:     Episodes,
	}
}

// ToLinks transforma datos de un servidor de video en un DTO LinkSource.
func (m *Maper) ToLinks(Server string, URL string, Code string) dto.LinkSource {
	return dto.LinkSource{
		Server: Server,
		URL:    URL,
		Code:   Code,
	}
}

// ToLinkEpisode transforma datos de enlaces de episodio en un DTO LinkResponse.
// Agrupa todos los enlaces de diferentes servidores para un episodio específico.
func (m *Maper) ToLinkEpisode(ID string, Title string, Episode uint, Links []dto.LinkSource) dto.LinkResponse {
	return dto.LinkResponse{
		ID:      ID,
		Episode: Episode,
		Title:   Title,
		Link:    Links,
	}
}

// ToRecentEpisode transforma datos de episodio reciente en un DTO EpisodeListResponse.
// Incluye información resumida para mostrar en listados de episodios recientes.
func (m *Maper) ToRecentEpisode(ID string, Title string, Chapter string, Episode int, Image string) dto.EpisodeListResponse {
	return dto.EpisodeListResponse{
		ID:      ID,
		Title:   Title,
		Chapter: Chapter,
		Episode: Episode,
		Image:   Image,
	}
}
