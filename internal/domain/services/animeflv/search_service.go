// Package animeflv - search_service.go
// Este archivo implementa el servicio de búsqueda de animes.
// Contiene la lógica de negocio para validar parámetros de búsqueda,
// normalizar entradas (convertir a minúsculas, manejar paginación),
// implementar caché distribuido y delegar al scraper para obtener los resultados.
package animeflv

import (
	"context"
	"fmt"
	"strings"

	"github.com/dst3v3n/api-anime/internal/domain/dto"
	"github.com/dst3v3n/api-anime/internal/ports"
)

// searchService encapsula la lógica de búsqueda de animes.
// Implementa caché distribuido (Valkey) para optimizar búsquedas frecuentes
// y reduce la carga al scraper mediante almacenamiento de resultados por página.
type searchService struct {
	scraper ports.ScraperPort
	cache   ports.CachePort
}

// SearchAnime realiza una búsqueda de animes con validaciones, transformaciones y caché.
// Valida que el nombre no esté vacío, normaliza el texto a minúsculas, maneja
// la paginación por defecto, intenta recuperar del caché y consulta al scraper si es necesario.
func (search *searchService) SearchAnime(anime *string, page *uint) (dto.AnimeResponse, error) {
	if *anime == "" {
		return dto.AnimeResponse{}, fmt.Errorf("el nombre del anime no puede estar vacío")
	}
	*anime = strings.ToLower(*anime)
	*anime = strings.ReplaceAll(*anime, " ", "-")

	if *page == 0 {
		*page = 1
	}
	pageStr := fmt.Sprintf("%d", *page)

	cacheKey := fmt.Sprintf("search-anime-%s-page-%d", *anime, *page)
	ctx := context.Background()

	var result dto.AnimeResponse

	if err := search.cache.Get(ctx, cacheKey, &result); err == nil {
		if len(result.Animes) > 0 {
			return result, nil
		}
	}

	result, err := search.scraper.SearchAnime(*anime, pageStr)
	if err != nil {
		return result, err
	}

	_ = search.cache.Set(ctx, cacheKey, result)

	return result, nil
}

// Search obtiene todos los animes disponibles sin filtros de búsqueda con caché.
// Intenta recuperar del caché primero, y si no está disponible, consulta al scraper
// y almacena el resultado en caché para futuras solicitudes.
func (search *searchService) Search() (dto.AnimeResponse, error) {
	catcheKey := "search-anime-all"
	ctx := context.Background()

	var result dto.AnimeResponse

	if err := search.cache.Get(ctx, catcheKey, &result); err == nil {
		if len(result.Animes) > 0 {
			return result, nil
		}
	}

	result, err := search.scraper.Search()
	if err != nil {
		return result, err
	}

	_ = search.cache.Set(ctx, catcheKey, result)

	return result, nil
}
