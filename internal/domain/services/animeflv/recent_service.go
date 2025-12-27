// Package animeflv - recent_service.go
// Este archivo implementa el servicio para obtener contenido reciente.
// Proporciona métodos para acceder a animes y episodios recientemente
// agregados al sitio, utilizando caché distribuido para optimizar consultas frecuentes
// y reducir la carga en el scraper.
package animeflv

import (
	"context"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/ports"
)

// recentService encapsula la lógica para obtener contenido reciente.
// Implementa caché distribuido (Valkey) para almacenar temporalmente resultados
// de animes y episodios recientes, mejorando el rendimiento de consultas repetidas.
type recentService struct {
	scraper ports.ScraperPort
	cache   ports.CachePort
}

// RecentAnime obtiene la lista de animes recientemente agregados con caché.
// Intenta recuperar del caché primero, y si no está disponible, consulta al scraper
// y almacena el resultado en caché para futuras solicitudes.
func (recent *recentService) RecentAnime() ([]dto.AnimeStruct, error) {
	cacheKey := "recent-anime"
	ctx := context.Background()

	var result []dto.AnimeStruct
	if err := recent.cache.Get(ctx, cacheKey, &result); err == nil {
		if len(result) > 0 {
			return result, nil
		}
	}

	result, err := recent.scraper.RecentAnime()
	if err != nil {
		return nil, err
	}

	_ = recent.cache.Set(ctx, cacheKey, result)

	return result, nil
}

// RecentEpisode obtiene la lista de episodios recientemente publicados con caché.
// Intenta recuperar del caché primero, y si no está disponible, consulta al scraper
// y almacena el resultado en caché para futuras solicitudes.
func (recent *recentService) RecentEpisode() ([]dto.EpisodeListResponse, error) {
	cacheKey := "recent-episode"
	ctx := context.Background()

	var result []dto.EpisodeListResponse
	if err := recent.cache.Get(ctx, cacheKey, &result); err == nil {
		if len(result) > 0 {
			return result, nil
		}
	}

	result, err := recent.scraper.RecentEpisode()
	if err != nil {
		return nil, err
	}

	_ = recent.cache.Set(ctx, cacheKey, result)

	return result, nil
}
