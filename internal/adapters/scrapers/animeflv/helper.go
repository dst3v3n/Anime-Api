// Package animeflv - helper.go
// Este archivo contiene funciones auxiliares de utilidad para procesar y transformar datos
// obtenidos del scraping de AnimeFlv. Incluye funciones para:
// - Parsear valores numéricos (episodios, puntuaciones)
// - Extraer IDs y números de episodio desde URLs
// - Construir URLs con parámetros de consulta
// - Manipular strings para limpiar y formatear datos
package animeflv

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

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

func extractEpisodeNumber(href string) (int, error) {
	parts := strings.Split(href, "-")
	if len(parts) > 1 {
		episodeStr := parts[len(parts)-1]
		episodeNum, err := strconv.Atoi(episodeStr)
		if err != nil {
			return 0, fmt.Errorf("invalid episode number format %q: %w", episodeStr, err)
		}
		return episodeNum, nil
	}
	return 0, fmt.Errorf("invalid href format for episode extraction: %s", href)
}

func removeTrailingNumber(id string) string {
	if lastDash := strings.LastIndex(id, "-"); lastDash != -1 {
		suffix := id[lastDash+1:]
		if _, err := strconv.Atoi(suffix); err == nil {
			return id[:lastDash]
		}
	}
	return id
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
