package animeflv

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/base"
	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"
	"github.com/gocolly/colly"
)

type ParserAnimeFlvStruct struct {
	base.ParserFuc
}

func getCategory(category string) dto.CategoryAnime {
	categoryMap := map[string]dto.CategoryAnime{
		"Anime":    dto.Anime,
		"Ova":      dto.Ova,
		"Película": dto.Pelicula,
	}

	return categoryMap[category]
}

func getStatus(status string) dto.StatusAnime {
	statusMap := map[string]dto.StatusAnime{
		"En emisión": dto.Emision,
		"Finalizado": dto.Finalizado,
	}
	return statusMap[status]
}

func replaceParenthesis(text string) string {
	text = strings.ReplaceAll(text, "(", "")
	text = strings.ReplaceAll(text, ")", "")
	return text
}

func getEpisodeInfo(scriptContent string) ([]int, string) {
	var episodios []int
	var nextEpisode string

	episodesRegex := regexp.MustCompile(`var episodes = (\[\[.*?\]\]);`)
	if matches := episodesRegex.FindStringSubmatch(scriptContent); len(matches) > 1 {
		var episodes [][]int
		if err := json.Unmarshal([]byte(matches[1]), &episodes); err == nil {
			for _, ep := range episodes {
				if len(ep) >= 2 {
					episodios = append(episodios, ep[0])
				}
			}
		}
	}

	animeInfoRegex := regexp.MustCompile(`var anime_info = (\[.*?\]);`)
	if matches := animeInfoRegex.FindStringSubmatch(scriptContent); len(matches) > 1 {
		var animeInfo []interface{}
		json.Unmarshal([]byte(matches[1]), &animeInfo)
		if len(animeInfo) >= 4 {
			nextEpisode = animeInfo[3].(string)
		}
	}

	return episodios, nextEpisode
}

func (ParseAnimeFlv ParserAnimeFlvStruct) ParserSearchAnime(url string) ([]dto.SearchAnimeResponse, error) {
	c := colly.NewCollector()
	var results []dto.SearchAnimeResponse

	c.OnHTML("ul.ListAnimes > li > article", func(e *colly.HTMLElement) {
		href := e.ChildAttr("a", "href")
		id := strings.Split(href, "/")[2]
		title := e.ChildText("h3.Title")
		category := e.ChildText("div.Description span.Type")
		category = cases.Title(language.Spanish).String(strings.ToLower(category))
		punctuation := e.ChildText("span.fa-star")
		puntuacion, _ := strconv.ParseFloat(punctuation, 64)
		image := e.ChildAttr("img", "src")
		sipnopsis := e.ChildText("div.Description p:nth-child(3)")
		results = append(results, dto.SearchAnimeResponse{
			ID:         id,
			Title:      title,
			Sipnopsis:  sipnopsis,
			Tipo:       getCategory(category),
			Puctuation: puntuacion,
			Image:      image,
		})
	})

	err := c.Visit(url)

	return results, err
}

func (ParseAnimeFlv ParserAnimeFlvStruct) ParserAnimeInfo(url string, idAnime string) (dto.ResponseAnimeInfo, error) {
	c := colly.NewCollector()

	// Variables para almacenar temporalmente los datos
	var title, category, sipnopsis, estado, image string
	var puntuacion float64
	var animeRelated []dto.AnimeRelated
	var generos []string
	var episodes []int
	var nextEpisode string
	var result dto.ResponseAnimeInfo
	var scrapeError error

	// Procesar scripts primero
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text
		if strings.Contains(scriptContent, "var episodes") {
			episodes, nextEpisode = getEpisodeInfo(scriptContent)
		}
	})

	// Procesar el contenido principal
	c.OnHTML("div.Body", func(e *colly.HTMLElement) {
		title = e.ChildText("h1.Title")
		category = e.ChildText("div.Container span.Type")
		sipnopsis = e.ChildText("div.Description p")
		estado = e.ChildText("span.fa-tv")
		punctuation := e.ChildText("span.vtprmd")
		puntuacion, _ = strconv.ParseFloat(punctuation, 64)
		image = e.ChildAttr("div.Image img", "src")

		e.ForEach("nav.Nvgnrs a", func(i int, el *colly.HTMLElement) {
			generos = append(generos, el.Text)
		})

		e.ForEach("ul.ListAnmRel > li", func(i int, el *colly.HTMLElement) {
			href := el.ChildAttr("a", "href")
			id := strings.Split(href, "/")[2]
			titleRel := el.ChildText("a")
			textoCompleto := el.Text
			categoryRel := strings.TrimSpace(strings.Replace(textoCompleto, titleRel, "", 1))
			categoryRel = replaceParenthesis(categoryRel)
			animeRelated = append(animeRelated, dto.AnimeRelated{
				ID:       id,
				Title:    titleRel,
				Category: categoryRel,
			})
		})
	})

	// Construir el resultado después de procesar todo
	c.OnScraped(func(r *colly.Response) {
		result = dto.ResponseAnimeInfo{
			SearchAnimeResponse: dto.SearchAnimeResponse{
				ID:         idAnime,
				Title:      title,
				Sipnopsis:  sipnopsis,
				Tipo:       getCategory(category),
				Puctuation: puntuacion,
				Image:      image,
			},
			AnimeRelated: animeRelated,
			Generos:      generos,
			Estado:       getStatus(estado),
			NextEpisode:  nextEpisode,
			Episodes:     episodes,
		}
	})

	// Manejo de errores
	c.OnError(func(r *colly.Response, error error) {
		if error != nil {
			scrapeError = error
		}
	})

	err := c.Visit(url)

	if scrapeError != nil {
		return dto.ResponseAnimeInfo{}, scrapeError
	}

	return result, err
}
