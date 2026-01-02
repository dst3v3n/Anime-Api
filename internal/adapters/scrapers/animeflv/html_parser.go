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
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dst3v3n/api-anime/internal/domain/dto"
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

// Parser es el componente principal de análisis HTML.
// Contiene un mapper para transformar datos extraídos en DTOs.
type Parser struct {
	mapper *Maper
}

// NewParser crea una nueva instancia del parser HTML.
// Inicializa el mapper interno para la transformación de datos.
func NewParser() *Parser {
	return &Parser{
		mapper: NewMaper(),
	}
}

// ParseAnime extrae información de animes desde HTML.
// Utilizado tanto para resultados de búsqueda como para animes recientes.
// Extrae: ID, título, sinopsis, tipo, puntuación e imagen de cada anime.
func (p *Parser) ParseAnime(htmlElement io.Reader) ([]dto.AnimeStruct, error) {
	result, err := p.ParseAnimeWithPagination(htmlElement)
	if err != nil {
		return result.Animes, fmt.Errorf("error al parsear animes: %w", err)
	}

	return result.Animes, nil
}

// ParseAnimeWithPagination extrae información de animes desde HTML con soporte de paginación.
// Utilizado tanto para resultados de búsqueda como para animes recientes.
// Proceso:
//  1. Parsea el HTML proporcionado en un documento goquery
//  2. Busca todos los artículos de anime usando selectores CSS
//  3. Para cada artículo extrae:
//     - ID desde el atributo href del enlace
//     - Título del anime
//     - Sinopsis descripción
//     - Tipo/Categoría (Anime, Película, OVA, Especial)
//     - Puntuación/Calificación
//     - URL de la imagen
//  4. Mapea cada elemento a AnimeStruct
//  5. Extrae información de paginación buscando el número de la penúltima página
//  6. Retorna error si no encuentra animes en el HTML
//
// Parámetros:
//   - htmlElement: reader del HTML con información de animes
//
// Retorna: AnimeResponse con lista de animes, número total de páginas, y error si falla
func (p *Parser) ParseAnimeWithPagination(htmlElement io.Reader) (dto.AnimeResponse, error) {
	results := dto.AnimeResponse{}
	doc, err := goquery.NewDocumentFromReader(htmlElement)
	if err != nil {
		return results, fmt.Errorf("error al crear documento desde HTML: %w", err)
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

		results.Animes = append(results.Animes, p.mapper.ToAnime(id, title, sinopsis, tipo, punctuation, image))
	})

	doc.Find("div.NvCnAnm ul.pagination").Each(func(_ int, s *goquery.Selection) {
		penultimoStr := s.Find("li:nth-last-child(2) a").Text()
		penultimo, err := parseUint(penultimoStr)
		results.TotalPages = penultimo
		if err != nil {
			return
		}
	})

	if len(results.Animes) == 0 {
		return results, fmt.Errorf("no se encontraron animes en el HTML proporcionado")
	}

	return results, nil
}

// ParseAnimeInfo extrae información completa de un anime específico.
// Procesa tanto el HTML visible como los scripts JavaScript embebidos para obtener:
// - Información básica (título, sinopsis, tipo, puntuación, imagen)
// - Géneros del anime
// - Estado de emisión y fecha del próximo episodio
// - Lista de episodios disponibles
// - Animes relacionados (secuelas, precuelas, spin-offs)
func (p *Parser) ParseAnimeInfo(htmlElement io.Reader, idAnime string) (dto.AnimeInfoResponse, error) {
	doc, err := goquery.NewDocumentFromReader(htmlElement)
	if err != nil {
		return dto.AnimeInfoResponse{}, fmt.Errorf("error al crear documento desde HTML: %w", err)
	}

	result := &ParseResult{}

	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		scriptContent := s.Text()

		if strings.Contains(scriptContent, "var episodes") {
			episodes, nextEpisode, err := episodeInfo(scriptContent)
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

	if len(resultFinal.Title) == 0 {
		return dto.AnimeInfoResponse{}, fmt.Errorf("no se pudo parsear la información del anime del HTML proporcionado")
	}

	return resultFinal, nil
}

// ParseLinks extrae los enlaces de reproducción de un episodio específico.
// Analiza scripts JavaScript embebidos en la página para obtener URLs de múltiples servidores
// de video (Zippyshare, Mega, Google Drive, etc.) junto con sus códigos de embed.
// Proceso:
//  1. Parsea el HTML de la página de reproducción
//  2. Busca etiquetas <script> que contienen la variable "var videos"
//  3. Utiliza scriptLinksEpisode() para extraer y parsear el JSON embebido
//  4. Extrae el título del anime del HTML
//  5. Valida que se hayan encontrado enlaces, retorna error si no hay
//  6. Mapea los resultados a LinkResponse
//
// Parámetros:
//   - htmlElement: reader del HTML de la página de reproducción del episodio
//   - idAnime: identificador del anime
//   - episodeNum: número del episodio
//
// Retorna: LinkResponse con título, número de episodio y lista de enlaces, o error si falla
func (p *Parser) ParseLinks(htmlElement io.Reader, idAnime string, episodeNum uint) (dto.LinkResponse, error) {
	doc, err := goquery.NewDocumentFromReader(htmlElement)
	if err != nil {
		return dto.LinkResponse{}, fmt.Errorf("error al crear documento desde HTML: %w", err)
	}

	result := &ParseEpisodeLinksResult{
		ID:      idAnime,
		Episode: episodeNum,
	}

	doc.Find("script[type=\"text/javascript\"]").Each(func(_ int, s *goquery.Selection) {
		scriptContent := s.Text()
		if strings.Contains(scriptContent, "var videos") {
			links, err := scriptLinksEpisode(scriptContent)
			if err != nil {
				return
			}
			result.links = links
		}
	})

	doc.Find(selectorBodyContainer).Each(func(_ int, s *goquery.Selection) {
		result.Title, _ = s.Find(selectorInfoTitle).Html()
	})

	if len(result.links) == 0 {
		return dto.LinkResponse{}, fmt.Errorf("no se pudo parsear los enlaces del episodio del HTML proporcionado")
	}

	return p.mapper.ToLinkEpisode(result.ID, result.Title, result.Episode, result.links), nil
}

// ParseRecentEpisode extrae la lista de episodios recientemente publicados desde la página principal.
// Obtiene información resumida de cada episodio para mostrar en listados recientes.
// Proceso:
//  1. Parsea el HTML de la página principal de AnimeFlv
//  2. Busca elementos <li> dentro del selector selectorEpisodeList
//  3. Para cada elemento, extrae:
//     - ID del anime desde el href (con normalización para remover número de episodio)
//     - Número de episodio desde el href
//     - Título del anime
//     - Designación del capítulo (ej: "Cap. 1050")
//     - URL de la imagen de portada
//  4. Utiliza el mapper para crear estructuras EpisodeListResponse
//  5. Retorna error si no encuentra ningún episodio reciente
//
// Parámetros:
//   - htmlElement: reader del HTML de la página principal
//
// Retorna: slice de EpisodeListResponse con información resumida de episodios recientes
func (p *Parser) ParseRecentEpisode(htmlElement io.Reader) ([]dto.EpisodeListResponse, error) {
	doc, err := goquery.NewDocumentFromReader(htmlElement)
	if err != nil {
		return []dto.EpisodeListResponse{}, fmt.Errorf("error al crear documento desde HTML: %w", err)
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

	if len(result) == 0 {
		return result, fmt.Errorf("no se encontraron episodios recientes en el HTML proporcionado")
	}
	return result, nil
}
