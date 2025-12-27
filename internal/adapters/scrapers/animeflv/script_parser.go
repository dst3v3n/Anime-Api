// Package animeflv - script_parser.go
// Este archivo se especializa en extraer y parsear datos embebidos en etiquetas <script>
// del HTML de AnimeFlv. Utiliza expresiones regulares para encontrar variables JavaScript
// que contienen información estructurada en formato JSON sobre:
// - Lista de episodios disponibles
// - Información de próximos episodios
// - Enlaces de servidores de video para reproducción
package animeflv

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
)

// scriptEpisodeList extrae la lista de episodios desde una variable JavaScript.
// Busca y parsea la variable "var episodes = [[...]]" que contiene un array bidimensional
// donde cada elemento tiene [episodeNumber, id]. Retorna solo los números de episodio.
func scriptEpisodeList(scriptContent string) ([]int, error) {
	var episodios []int
	episodesRegex := regexp.MustCompile(`var episodes = (\[\[.*?\]\]);`)
	if matches := episodesRegex.FindStringSubmatch(scriptContent); len(matches) > 1 {
		var episodes [][]int
		if err := json.Unmarshal([]byte(matches[1]), &episodes); err != nil {
			return nil, fmt.Errorf("error al parsear JSON de episodios: %w", err)
		}
		for _, ep := range episodes {
			if len(ep) >= 2 {
				episodios = append(episodios, ep[0])
			}
		}
	}
	return episodios, nil
}

// scriptInfo extrae información adicional del anime desde una variable JavaScript.
// Busca y parsea la variable "var anime_info = [...]" que contiene un array con
// datos del anime. El cuarto elemento (índice 3) contiene la fecha del próximo episodio.
func scriptInfo(scriptContent string) (string, error) {
	var nextEpisode string
	animeInfoRegex := regexp.MustCompile(`var anime_info = (\[.*?\]);`)
	if matches := animeInfoRegex.FindStringSubmatch(scriptContent); len(matches) > 1 {
		var animeInfo []any
		if err := json.Unmarshal([]byte(matches[1]), &animeInfo); err != nil {
			return "", fmt.Errorf("error al parsear JSON de información del anime: %w", err)
		}
		if len(animeInfo) >= 4 {
			nextEpisode = animeInfo[3].(string)
		}
	}
	return nextEpisode, nil
}

// scriptLinksEpisode extrae los enlaces de video desde una variable JavaScript.
// Busca y parsea la variable "var videos = {...}" que contiene un objeto JSON con
// servidores de video (SUB, LAT, etc.). Cada servidor tiene información de URL,
// código de embed y otras propiedades. Retorna una lista de fuentes de enlaces.
func scriptLinksEpisode(scriptContent string) ([]dto.LinkSource, error) {
	toLinkSource := []dto.LinkSource{}
	linkEpisodeRegex := regexp.MustCompile(`var videos = (\{.*?\});`)
	if matches := linkEpisodeRegex.FindStringSubmatch(scriptContent); len(matches) > 1 {
		var videos Videos
		if err := json.Unmarshal([]byte(matches[1]), &videos); err != nil {
			return nil, fmt.Errorf("error al parsear JSON de enlaces de video: %w", err)
		}
		for _, linkVideo := range videos.SUB {
			toLinkSource = append(toLinkSource, dto.LinkSource{
				Server: linkVideo.Server,
				URL:    linkVideo.URL,
				Code:   linkVideo.Code,
			})
		}
	}
	return toLinkSource, nil
}
