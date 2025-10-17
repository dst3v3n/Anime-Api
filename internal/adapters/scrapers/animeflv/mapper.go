package animeflv

import "github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"

type Maper struct{}

func NewMaper() *Maper {
	return &Maper{}
}

func (m *Maper) ToSearchAnime(ID string, Title string, Sipnopsis string, Tipo string, Puctuation float64, Image string) dto.SearchAnimeResponse {
	return dto.SearchAnimeResponse{
		ID:         ID,
		Title:      Title,
		Sipnopsis:  Sipnopsis,
		Tipo:       dto.CategoryAnime(Tipo),
		Puctuation: Puctuation,
		Image:      Image,
	}
}

func (m *Maper) ToAnimeInfo(ID string, Title string, Sipnopsis string, Tipo string, Puctuation float64, Image string, AnimeRelated []dto.AnimeRelated, Generos []string, Estado string, Episodes []int, NextEpisode string) dto.ResponseAnimeInfo {
	return dto.ResponseAnimeInfo{
		SearchAnimeResponse: dto.SearchAnimeResponse{
			ID:         ID,
			Title:      Title,
			Sipnopsis:  Sipnopsis,
			Tipo:       dto.CategoryAnime(Tipo),
			Puctuation: Puctuation,
			Image:      Image,
		},
		AnimeRelated: AnimeRelated,
		Generos:      Generos,
		Estado:       dto.StatusAnime(Estado),
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

func (m *Maper) ToLinkEpisode(ID string, Title string, Episode int, Links []dto.LinkSource) dto.GetLinkResponse {
	return dto.GetLinkResponse{
		ID:      ID,
		Episode: Episode,
		Title:   Title,
		Link:    Links,
	}
}
