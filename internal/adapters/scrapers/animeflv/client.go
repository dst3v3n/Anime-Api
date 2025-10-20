// Package animeflv implementa un cliente scraper para el sitio web AnimeFlv.
// Este archivo (client.go) contiene la estructura principal del cliente y la implementación
// de todos los métodos definidos en el port ScraperPort. Se encarga de realizar las
// peticiones HTTP a AnimeFlv y delegar el parsing del HTML al componente Parser.
package animeflv

import (
	"fmt"
	"net/http"

	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"
	"github.com/dst3v3n/api-anime/internal/ports"
)

type Client struct {
	config Config
	parser *Parser
}

func NewClient() ports.ScraperPort {
	return &Client{
		config: Config{
			BaseURL:       "https://www3.animeflv.net",
			SearchURL:     "https://www3.animeflv.net/browse",
			AnimeInfoURL:  "https://www3.animeflv.net/anime",
			VerEpisodeURL: "https://www3.animeflv.net/ver",
		},
		parser: NewParser(),
	}
}

func (c *Client) SearchAnime(anime string, page string) ([]dto.AnimeResponse, error) {
	if page == "" {
		page = "1"
	}

	params := map[string]string{
		"page": page,
		"q":    anime,
	}

	searchURL := buildURL(c.config.SearchURL, params)
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.parser.ParseAnime(resp.Body)
}

func (c *Client) AnimeInfo(idAnime string) (dto.AnimeInfoResponse, error) {
	searhURL := c.config.AnimeInfoURL + "/" + idAnime
	resp, err := http.Get(searhURL)
	if err != nil {
		return dto.AnimeInfoResponse{}, err
	}
	defer resp.Body.Close()

	return c.parser.ParseAnimeInfo(resp.Body, idAnime)
}

func (c *Client) Links(idAnime string, episode int) (dto.LinkResponse, error) {
	searchURL := fmt.Sprintf("%s/%s-%d", c.config.VerEpisodeURL, idAnime, episode)
	resp, err := http.Get(searchURL)
	if err != nil {
		return dto.LinkResponse{}, err
	}

	defer resp.Body.Close()

	return c.parser.ParseLinks(resp.Body, idAnime, episode)
}

func (c *Client) RecentAnime() ([]dto.AnimeResponse, error) {
	searchURL := c.config.BaseURL
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return c.parser.ParseAnime(resp.Body)
}

func (c *Client) RecentEpisode() ([]dto.EpisodeListResponse, error) {
	searchURL := c.config.BaseURL
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.parser.ParseRecentEpisode(resp.Body)
}
