// Package animeflv - models.go
// Este archivo define las estructuras de datos internas utilizadas específicamente
// por el scraper de AnimeFlv. Incluye:
// - Config: Configuración de URLs del sitio
// - ParseResult: Estructura temporal para almacenar datos durante el parsing de información de anime
// - ParseEpisodeLinksResult: Estructura temporal para almacenar enlaces de episodios
// - VideoServer y Videos: Estructuras para deserializar JSON embebido en scripts del sitio
package animeflv

import "github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"

type Config struct {
	BaseURL       string
	SearchURL     string
	AnimeInfoURL  string
	VerEpisodeURL string
}

type ParseResult struct {
	title, category, sipnopsis, status, image string
	punctuacion                               float64
	animeRelated                              []dto.AnimeRelated
	genres                                    []string
	episodes                                  []int
	nextEpisode                               string
}

type ParseEpisodeLinksResult struct {
	ID      string
	Title   string
	Episode int
	links   []dto.LinkSource
}

type VideoServer struct {
	Server      string `json:"server"`
	Title       string `json:"title"`
	Ads         int    `json:"ads"`
	URL         string `json:"url,omitempty"`
	AllowMobile bool   `json:"allow_mobile"`
	Code        string `json:"code,omitempty"`
}

type Videos struct {
	SUB []VideoServer `json:"SUB"`
}
