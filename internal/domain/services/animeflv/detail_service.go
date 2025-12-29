// Package animeflv - detail_service.go
// Este archivo implementa el servicio de detalles de anime.
// Contiene la lógica de negocio para obtener información detallada de animes
// y enlaces de episodios, aplicando validaciones de parámetros, caché distribuido,
// y normalizaciones antes de delegar al scraper.
package animeflv

import (
	"context"
	"fmt"
	"strings"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/ports"
)

// detailService encapsula la lógica para obtener detalles de anime y episodios.
// Utiliza caché distribuido (Valkey) para optimizar consultas recurrentes
// y reduce la carga al scraper mediante almacenamiento temporal de resultados.
type detailService struct {
	scraper     ports.ScraperPort
	cache       ports.CachePort
	enableCache bool
}

// AnimeInfo obtiene información completa de un anime aplicando validaciones y caché.
// Verifica que el ID no esté vacío, lo normaliza a minúsculas, intenta recuperar
// del caché y, si no existe, consulta al scraper y almacena el resultado en caché.
func (detail *detailService) AnimeInfo(ctx context.Context, idAnime string) (dto.AnimeInfoResponse, error) {
	if idAnime == "" {
		return dto.AnimeInfoResponse{}, fmt.Errorf("el ID del anime no puede estar vacío")
	}

	id := strings.ToLower(strings.TrimSpace(idAnime))
	cacheKey := fmt.Sprintf("anime-info-%s", id)

	if detail.enableCache {
		var result dto.AnimeInfoResponse
		if err := detail.cache.Get(ctx, cacheKey, &result); err == nil {
			if len(result.ID) > 0 {
				return result, nil
			}
		}
	}

	result, err := detail.scraper.AnimeInfo(ctx, id)
	if err != nil {
		return dto.AnimeInfoResponse{}, err
	}

	if detail.enableCache {
		_ = detail.cache.Set(ctx, cacheKey, result)
	}

	return result, nil
}

// Links obtiene los enlaces de reproducción de un episodio específico con caché.
// Valida que el ID del anime no esté vacío, normaliza a minúsculas, intenta recuperar
// del caché y, si no existe, consulta al scraper y almacena el resultado en caché
// para futuras solicitudes del mismo episodio.
func (detail *detailService) Links(ctx context.Context, idAnime string, episode uint) (dto.LinkResponse, error) {
	if idAnime == "" {
		return dto.LinkResponse{}, fmt.Errorf("el ID del anime no puede estar vacío")
	}

	id := strings.ToLower(strings.TrimSpace(idAnime))
	cacheKey := fmt.Sprintf("links-%s-%d", id, episode)

	if detail.enableCache {
		var result dto.LinkResponse
		if err := detail.cache.Get(ctx, cacheKey, &result); err == nil {
			if len(result.ID) > 0 {
				return result, nil
			}
		}
	}

	result, err := detail.scraper.Links(ctx, id, episode)
	if err != nil {
		return dto.LinkResponse{}, err
	}

	if detail.enableCache {
		_ = detail.cache.Set(ctx, cacheKey, result)
	}

	return result, nil
}
