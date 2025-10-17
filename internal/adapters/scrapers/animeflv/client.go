package animeflv

import (
	"fmt"
	"net/url"
	"strings"

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

func (c *Client) SearchAnime(anime string, page string) ([]dto.SearchAnimeResponse, error) {
	if page == "" {
		page = "1"
	}

	params := map[string]string{
		"page": page,
		"q":    anime,
	}

	searchURL := buildURL(c.config.SearchURL, params)
	return c.parser.ParseSearchAnime(searchURL)
}

func (c *Client) AnimeInfo(idAnime string) (dto.ResponseAnimeInfo, error) {
	idAnime = strings.ToLower(idAnime)
	animeURL := c.config.AnimeInfoURL + "/" + idAnime
	return c.parser.ParseAnimeInfo(animeURL, idAnime)
}

func (c *Client) GetLinks(idAnime string, episode int) (dto.GetLinkResponse, error) {
	idAnime = strings.ToLower(idAnime)
	episodeURL := fmt.Sprintf("%s/%s-%d", c.config.VerEpisodeURL, idAnime, episode)
	return c.parser.ParseGetLinks(episodeURL, idAnime, episode)
}

func buildURL(baseURL string, params map[string]string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	query := u.Query()
	for key, value := range params {
		query.Add(key, value)
	}

	u.RawQuery = query.Encode()
	return u.String()
}
