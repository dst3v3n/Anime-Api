// Package animeflv provides HTML parsing functionality for the AnimeFlv website.
package animeflv

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"
	"github.com/gocolly/colly"
)

// CSS selectors constants
const (
	selectorSearchArticle      = "ul.ListAnimes > li > article"
	selectorArticleLink        = "a"
	selectorArticleTitle       = "h3.Title"
	selectorArticleCategory    = "div.Description span.Type"
	selectorArticlePunctuation = "span.fa-star"
	selectorArticleImage       = "img"
	selectorArticleSynopsis    = "div.Description p:nth-child(3)"

	selectorBodyContainer   = "div.Body"
	selectorInfoTitle       = "h1.Title"
	selectorInfoCategory    = "div.Container span.Type"
	selectorInfoSynopsis    = "div.Description p"
	selectorInfoStatus      = "span.fa-tv"
	selectorInfoPunctuation = "span.vtprmd"
	selectorInfoImage       = "div.Image img"
	selectorInfoGenres      = "nav.Nvgnrs a"
	selectorInfoRelated     = "ul.ListAnmRel > li"
)

// Parser maneja el parseo de HTML de AnimeFlv
type Parser struct {
	collector *colly.Collector
	mapper    *Maper
}

// NewParser crea una nueva instancia del parser
func NewParser() *Parser {
	return &Parser{
		collector: colly.NewCollector(
			colly.AllowedDomains("www3.animeflv.net")),
		mapper: NewMaper(),
	}
}

// ParseSearchAnime parsea los resultados de búsqueda
func (p *Parser) ParseSearchAnime(url string) ([]dto.SearchAnimeResponse, error) {
	c := p.collector.Clone()
	var results []dto.SearchAnimeResponse
	var parseErrors []error

	c.OnHTML(selectorSearchArticle, func(e *colly.HTMLElement) {
		// Extract ID
		href := e.ChildAttr(selectorArticleLink, "href")
		id, err := extractID(href)
		if err != nil {
			parseErrors = append(parseErrors, fmt.Errorf("failed to extract ID from href %s: %w", href, err))
			return
		}

		// Extract and normalize category
		title := e.ChildText(selectorArticleTitle)
		category := e.ChildText(selectorArticleCategory)
		category = cases.Title(language.Spanish).String(strings.ToLower(category))

		// Extract punctuation with fallback
		punctuation := e.ChildText(selectorArticlePunctuation)
		puntuacion, err := parseFloat(punctuation)
		if err != nil {
			puntuacion = 0.0 // Fallback to 0 if parsing fails
		}

		// Extract remaining fields
		image := e.ChildAttr(selectorArticleImage, "src")
		sipnopsis := e.ChildText(selectorArticleSynopsis)

		results = append(results, p.mapper.ToSearchAnime(id, title, sipnopsis, category, puntuacion, image))
	})

	if err := c.Visit(url); err != nil {
		return nil, fmt.Errorf("failed to visit URL %s: %w", url, err)
	}

	// If we have parse errors but also results, log them but don't fail
	if len(parseErrors) > 0 && len(results) == 0 {
		return nil, fmt.Errorf("all items failed to parse: %d errors", len(parseErrors))
	}

	return results, nil
}

// ParseAnimeInfo parsea la información detallada de un anime
func (p *Parser) ParseAnimeInfo(url string, idAnime string) (dto.ResponseAnimeInfo, error) {
	c := p.collector.Clone()
	result := &ParseResult{}
	var scrapeError error

	// Process scripts first
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text
		if strings.Contains(scriptContent, "var episodes") {
			episodes, nextEp, err := getEpisodeInfo(scriptContent)
			if err != nil {
				// Don't fail, just log and continue
				result.episodes = []int{}
				result.nextEpisode = ""
				return
			}
			result.episodes = episodes
			result.nextEpisode = nextEp
		}
	})

	// Process main content
	c.OnHTML(selectorBodyContainer, func(e *colly.HTMLElement) {
		result.title = e.ChildText(selectorInfoTitle)
		result.category = e.ChildText(selectorInfoCategory)
		result.sipnopsis = e.ChildText(selectorInfoSynopsis)
		result.estado = e.ChildText(selectorInfoStatus)

		punctuation := e.ChildText(selectorInfoPunctuation)
		puntuacion, err := parseFloat(punctuation)
		if err != nil {
			puntuacion = 0.0
		}
		result.puntuacion = puntuacion
		result.image = e.ChildAttr(selectorInfoImage, "src")

		// Extract genres
		e.ForEach(selectorInfoGenres, func(i int, el *colly.HTMLElement) {
			result.generos = append(result.generos, el.Text)
		})

		// Extract related anime
		e.ForEach(selectorInfoRelated, func(i int, el *colly.HTMLElement) {
			href := el.ChildAttr(selectorArticleLink, "href")
			id, err := extractID(href)
			if err != nil {
				return // Skip this item if ID extraction fails
			}

			titleRel := el.ChildText(selectorArticleLink)
			textoCompleto := el.Text
			categoryRel := strings.TrimSpace(strings.TrimPrefix(textoCompleto, titleRel))
			categoryRel = replaceParenthesis(categoryRel)

			result.animeRelated = append(result.animeRelated, dto.AnimeRelated{
				ID:       id,
				Title:    titleRel,
				Category: categoryRel,
			})
		})
	})

	// Build final result after processing
	var finalResultado dto.ResponseAnimeInfo
	c.OnScraped(func(r *colly.Response) {
		finalResultado = p.mapper.ToAnimeInfo(
			idAnime,
			result.title,
			result.sipnopsis,
			result.category,
			result.puntuacion,
			result.image,
			result.animeRelated,
			result.generos,
			result.estado,
			result.episodes,
			result.nextEpisode,
		)
	})

	// Error handling
	c.OnError(func(r *colly.Response, err error) {
		scrapeError = err
	})

	if err := c.Visit(url); err != nil {
		return dto.ResponseAnimeInfo{}, fmt.Errorf("failed to visit URL %s: %w", url, err)
	}

	if scrapeError != nil {
		return dto.ResponseAnimeInfo{}, fmt.Errorf("scraping error: %w", scrapeError)
	}

	return finalResultado, nil
}

func (p *Parser) ParseGetLinks(url string, idAnime string, episode int) (dto.GetLinkResponse, error) {
	c := p.collector.Clone()
	result := &ParseEpisodeLinksResult{
		ID:      idAnime,
		Episode: episode,
	}

	c.OnHTML("script[type=\"text/javascript\"]", func(e *colly.HTMLElement) {
		scriptContent := e.Text
		if strings.Contains(scriptContent, "var videos") {
			links, error := GetScriptLinksEpisode(scriptContent)
			if error != nil {
				return
			}
			result.links = links
		}
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		result.Title = e.ChildText(selectorInfoTitle)
	})

	var resultFinal dto.GetLinkResponse
	c.OnScraped(func(r *colly.Response) {
		resultFinal = p.mapper.ToLinkEpisode(result.ID, result.Title, result.Episode, result.links)
	})

	if err := c.Visit(url); err != nil {
		return dto.GetLinkResponse{}, fmt.Errorf("failed to visit URL %s: %w", url, err)
	}

	return resultFinal, nil
}

// Helper functions

func replaceParenthesis(text string) string {
	text = strings.ReplaceAll(text, "(", "")
	text = strings.ReplaceAll(text, ")", "")
	return text
}

func getEpisodeInfo(scriptContent string) ([]int, string, error) {
	episodios, err := GetScriptEpisodeList(scriptContent)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get episode list: %w", err)
	}

	nextEpisode, err := GetScriptInfo(scriptContent)
	if err != nil {
		return episodios, "", fmt.Errorf("failed to get next episode info: %w", err)
	}

	return episodios, nextEpisode, nil
}

func parseFloat(value string) (float64, error) {
	if value == "" {
		return 0.0, fmt.Errorf("empty string, cannot parse to float")
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0, fmt.Errorf("invalid float format %q: %w", value, err)
	}
	return parsed, nil
}

func extractID(href string) (string, error) {
	parts := strings.Split(href, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid href format: %s", href)
	}
	return parts[2], nil
}
