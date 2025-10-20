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
	"regexp"

	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"
)

func GetScriptEpisodeList(scriptContent string) ([]int, error) {
	var episodios []int

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
	return episodios, nil
}

func GetScriptInfo(scriptContent string) (string, error) {
	var nextEpisode string
	animeInfoRegex := regexp.MustCompile(`var anime_info = (\[.*?\]);`)
	if matches := animeInfoRegex.FindStringSubmatch(scriptContent); len(matches) > 1 {
		var animeInfo []any
		error := json.Unmarshal([]byte(matches[1]), &animeInfo)
		if error != nil {
			return "", error
		}
		if len(animeInfo) >= 4 {
			nextEpisode = animeInfo[3].(string)
		}
	}

	return nextEpisode, nil
}

func GetScriptLinksEpisode(scriptContent string) ([]dto.LinkSource, error) {
	toLinkSource := []dto.LinkSource{}
	linkEpisodeRegex := regexp.MustCompile(`var videos = (\{.*?\});`)
	if matches := linkEpisodeRegex.FindStringSubmatch(scriptContent); len(matches) > 1 {
		var videos Videos
		error := json.Unmarshal([]byte(matches[1]), &videos)
		if error != nil {
			return nil, error
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
