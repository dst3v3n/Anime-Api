package animeflv

import "github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"

type Config struct {
	BaseURL       string
	SearchURL     string
	AnimeInfoURL  string
	VerEpisodeURL string
}

type ParseResult struct {
	title, category, sipnopsis, estado, image string
	puntuacion                                float64
	animeRelated                              []dto.AnimeRelated
	generos                                   []string
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
