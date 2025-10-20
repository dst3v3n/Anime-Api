// Package animeflv - mapper.go
// Este archivo implementa el componente Mapper que se encarga de transformar
// datos crudos extraídos del HTML en estructuras DTO (Data Transfer Objects) bien definidas.
// Convierte los datos primitivos en objetos de dominio utilizables por la aplicación,
// proporcionando una capa de abstracción entre el scraping y la lógica de negocio.
package animeflv

import "github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"

type Maper struct{}

func NewMaper() *Maper {
	return &Maper{}
}

func (m *Maper) ToAnime(ID string, Title string, Sinopsis string, Tipo string, Punctuation float64, Image string) dto.AnimeResponse {
	return dto.AnimeResponse{
		ID:          ID,
		Title:       Title,
		Sinopsis:    Sinopsis,
		Type:        dto.CategoryAnime(Tipo),
		Punctuation: Punctuation,
		Image:       Image,
	}
}

func (m *Maper) ToAnimeInfo(ID string, Title string, Sinopsis string, Tipo string, Punctuation float64, Image string, AnimeRelated []dto.AnimeRelated, Generos []string, Estado string, Episodes []int, NextEpisode string) dto.AnimeInfoResponse {
	return dto.AnimeInfoResponse{
		AnimeResponse: dto.AnimeResponse{
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

func (m *Maper) ToLinks(Server string, URL string, Code string) dto.LinkSource {
	return dto.LinkSource{
		Server: Server,
		URL:    URL,
		Code:   Code,
	}
}

func (m *Maper) ToLinkEpisode(ID string, Title string, Episode int, Links []dto.LinkSource) dto.LinkResponse {
	return dto.LinkResponse{
		ID:      ID,
		Episode: Episode,
		Title:   Title,
		Link:    Links,
	}
}

func (m *Maper) ToRecentEpisode(ID string, Title string, Chapter string, Episode int, Image string) dto.EpisodeListResponse {
	return dto.EpisodeListResponse{
		ID:      ID,
		Title:   Title,
		Chapter: Chapter,
		Episode: Episode,
		Image:   Image,
	}
}
