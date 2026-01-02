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

// episodeInfo extrae información de episodios desde el contenido de un script JavaScript.
// Combina los resultados de scriptEpisodeList y scriptInfo para obtener
// tanto la lista de episodios disponibles como la fecha del próximo episodio.
func episodeInfo(scriptContent string) ([]int, string, error) {
	episodios, err := scriptEpisodeList(scriptContent)
	if err != nil {
		return nil, "", fmt.Errorf("error al obtener lista de episodios: %w", err)
	}
	nextEpisode, err := scriptInfo(scriptContent)
	if err != nil {
		return episodios, "", fmt.Errorf("error al obtener información del próximo episodio: %w", err)
	}
	return episodios, nextEpisode, nil
}

// parseFloat convierte una cadena a float64 con validación.
// Retorna error si la cadena está vacía o no tiene un formato numérico válido.
func parseFloat(value string) (float64, error) {
	if value == "" {
		return 0.0, fmt.Errorf("cadena vacía, no se puede convertir a float")
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0, fmt.Errorf("formato de float inválido %q: %w", value, err)
	}
	return parsed, nil
}

// extractID extrae el ID del anime desde una URL de AnimeFlv.
// Ejemplo: "/anime/one-piece-tv" -> "one-piece-tv"
func extractID(href string) (string, error) {
	parts := strings.Split(href, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("formato de href inválido: %s", href)
	}
	return parts[2], nil
}

// extractEpisodeNumber extrae el número de episodio desde una URL.
// Ejemplo: "/ver/one-piece-tv-1150" -> 1150
func extractEpisodeNumber(href string) (int, error) {
	parts := strings.Split(href, "-")
	if len(parts) > 1 {
		episodeStr := parts[len(parts)-1]
		episodeNum, err := strconv.Atoi(episodeStr)
		if err != nil {
			return 0, fmt.Errorf("formato de número de episodio inválido %q: %w", episodeStr, err)
		}
		return episodeNum, nil
	}
	return 0, fmt.Errorf("formato de href inválido para extracción de episodio: %s", href)
}

// removeTrailingNumber elimina el número final de un ID si existe.
// Ejemplo: "one-piece-tv-1150" -> "one-piece-tv"
// Útil para normalizar IDs de episodios a IDs de anime.
func removeTrailingNumber(id string) string {
	if lastDash := strings.LastIndex(id, "-"); lastDash != -1 {
		suffix := id[lastDash+1:]
		if _, err := strconv.Atoi(suffix); err == nil {
			return id[:lastDash]
		}
	}
	return id
}

// buildURL construye una URL completa agregando parámetros de consulta.
// Ejemplo: buildURL("https://example.com/search", {"q": "naruto", "page": "1"})
// retorna "https://example.com/search?q=naruto&page=1"
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

// parseUint convierte una cadena a uint con validación.
// parseUint convierte una cadena a uint con validación.
// Valida que la cadena tenga un formato numérico válido y sea un número no-negativo.
// Si la conversión es exitosa, retorna el valor como uint; de lo contrario, retorna error.
// Parámetros:
//   - value: cadena a convertir (ej: "5", "100", "99999")
//
// Retorna:
//   - uint: valor numérico sin signo convertido
//   - error: error si la cadena está vacía, contiene caracteres no numéricos, o es un número negativo
func parseUint(value string) (uint, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("formato de uint inválido %q: %w", value, err)
	}
	return uint(parsed), nil
}
