// Package animeflv - html_parser.go
// Este archivo contiene la lógica principal de parsing HTML del sitio AnimeFlv.
// Utiliza goquery para analizar el DOM HTML y extraer información estructurada sobre:
// - Búsqueda de animes
// - Información detallada de anime (géneros, estado, episodios, animes relacionados)
// - Enlaces de reproducción de episodios
// - Listado de episodios recientes
// Define todos los selectores CSS necesarios y coordina el proceso de extracción y mapeo de datos.
package animeflv

import (
	"html"
	"io"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"
)

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

	selectorEpisodeList        = "ul.ListEpisodios > li"
	selectorEpisodeListTitle   = "strong.Title"
	selectorEpisodeListChapter = "span.Capi"
)

type Parser struct {
	mapper *Maper
}

func NewParser() *Parser {
	return &Parser{
		mapper: NewMaper(),
	}
}

func (p *Parser) ParseAnime(htmlElement io.Reader) ([]dto.AnimeResponse, error) {
	results := []dto.AnimeResponse{}
	doc, err := goquery.NewDocumentFromReader(htmlElement)
	if err != nil {
		return results, err
	}

	doc.Find(selectorSearchArticle).Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Find(selectorArticleLink).Attr("href")
		id, err := extractID(href)
		if err != nil {
			return
		}

		image, _ := s.Find(selectorArticleImage).Attr("src")
		tipo, _ := s.Find(selectorArticleCategory).Html()
		title, _ := s.Find(selectorArticleTitle).Html()
		punctuationStr, _ := s.Find(selectorArticlePunctuation).Html()
		punctuation, err := parseFloat(punctuationStr)
		if err != nil {
			return
		}

		sinopsis, _ := s.Find(selectorArticleSynopsis).Html()
		sinopsis = html.UnescapeString(sinopsis)

		results = append(results, p.mapper.ToAnime(id, title, sinopsis, tipo, punctuation, image))
	})
	return results, nil
}

func (p *Parser) ParseAnimeInfo(htmlElement io.Reader, idAnime string) (dto.AnimeInfoResponse, error) {
	doc, err := goquery.NewDocumentFromReader(htmlElement)
	if err != nil {
		return dto.AnimeInfoResponse{}, err
	}

	result := &ParseResult{}

	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		scriptContent := s.Text()
		if strings.Contains(scriptContent, "var episodes") {
			episodes, nextEpisode, err := getEpisodeInfo(scriptContent)
			if err != nil {
				result.episodes = []int{}
				result.nextEpisode = ""
				return
			}
			result.episodes = episodes
			result.nextEpisode = nextEpisode
		}
	})

	doc.Find(selectorBodyContainer).Each(func(_ int, s *goquery.Selection) {
		result.title, _ = s.Find(selectorInfoTitle).Html()
		result.category, _ = s.Find(selectorInfoCategory).Html()
		result.image, _ = s.Find(selectorInfoImage).Attr("src")

		s.Find(selectorInfoGenres).Each(func(_ int, genreSel *goquery.Selection) {
			result.genres = append(result.genres, genreSel.Text())
		})

		sinopsis, _ := s.Find(selectorInfoSynopsis).Html()
		result.sipnopsis = html.UnescapeString(sinopsis)
		result.status, _ = s.Find(selectorInfoStatus).Html()
		punctuationStr, _ := s.Find(selectorInfoPunctuation).Html()
		result.punctuacion, err = parseFloat(punctuationStr)
		if err != nil {
			return
		}

		s.Find(selectorInfoRelated).Each(func(_ int, relatedSel *goquery.Selection) {
			href, _ := relatedSel.Find(selectorArticleLink).Attr("href")
			title := relatedSel.Find(selectorArticleLink).Text()

			id, err := extractID(href)
			if err != nil {
				return
			}

			fullText := relatedSel.Text()

			re := regexp.MustCompile(`\((.*?)\)`)
			matches := re.FindStringSubmatch(fullText)

			relationType := ""
			if len(matches) > 1 {
				relationType = matches[1]
			}

			result.animeRelated = append(result.animeRelated, dto.AnimeRelated{
				ID:       id,
				Title:    title,
				Category: relationType,
			})
		})
	})

	resultFinal := p.mapper.ToAnimeInfo(
		idAnime,
		result.title,
		result.sipnopsis,
		result.category,
		result.punctuacion,
		result.image,
		result.animeRelated,
		result.genres,
		result.status,
		result.episodes,
		result.nextEpisode,
	)

	return resultFinal, nil
}

func (p *Parser) ParseLinks(htmlElement io.Reader, idAnime string, episodeNum int) (dto.LinkResponse, error) {
	doc, err := goquery.NewDocumentFromReader(htmlElement)
	if err != nil {
		return dto.LinkResponse{}, err
	}

	result := &ParseEpisodeLinksResult{
		ID:      idAnime,
		Episode: episodeNum,
	}

	doc.Find("script[type=\"text/javascript\"]").Each(func(_ int, s *goquery.Selection) {
		scriptContent := s.Text()
		if strings.Contains(scriptContent, "var videos") {
			links, err := GetScriptLinksEpisode(scriptContent)
			if err != nil {
				return
			}
			result.links = links
		}
	})

	doc.Find(selectorBodyContainer).Each(func(_ int, s *goquery.Selection) {
		result.Title, _ = s.Find(selectorInfoTitle).Html()
	})

	return p.mapper.ToLinkEpisode(result.ID, result.Title, result.Episode, result.links), nil
}

func (p *Parser) ParseRecentEpisode(htmlElement io.Reader) ([]dto.EpisodeListResponse, error) {
	doc, err := goquery.NewDocumentFromReader(htmlElement)
	if err != nil {
		return []dto.EpisodeListResponse{}, err
	}
	result := []dto.EpisodeListResponse{}

	doc.Find(selectorEpisodeList).Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Find(selectorArticleLink).Attr("href")
		id, err := extractID(href)
		if err != nil {
			return
		}

		episode, err := extractEpisodeNumber(href)
		if err != nil {
			return
		}

		id = removeTrailingNumber(id)
		title, _ := s.Find(selectorEpisodeListTitle).Html()
		chapter := s.Find(selectorEpisodeListChapter).Text()
		image, _ := s.Find(selectorArticleImage).Attr("src")

		result = append(result, p.mapper.ToRecentEpisode(id, title, chapter, episode, image))
	})
	return result, nil
}
