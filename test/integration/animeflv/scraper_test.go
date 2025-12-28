// Package animeflv contiene tests de integración para el scraper de AnimeFlv.
// Este archivo (scraper_test.go) implementa pruebas que verifican la interacción
// del scraper con el sitio web real de AnimeFlv, validando búsqueda, obtención de
// información de anime, enlaces de episodios y listados de contenido reciente.
package animeflv

import (
	"context"
	"testing"

	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/animeflv"
)

func TestSearchAnimeScraper(t *testing.T) {
	testCases := []struct {
		name        string
		animeSearch string
		page        string
		wantError   bool
		description string
	}{
		{
			name:        "Prueba búsqueda de anime 'naruto' página 1",
			animeSearch: "naruto",
			page:        "1",
			wantError:   false,
			description: "debe buscar correctamente el anime 'naruto' en la página 1 sin errores",
		},
		{
			name:        "Prueba búsqueda de anime 'one piece' página 0",
			animeSearch: "one piece",
			page:        "0",
			wantError:   false,
			description: "debe buscar correctamente el anime 'naruto' en la página 0 sin errores",
		},
		{
			name:        "Prueba búsqueda de anime 'naruto hakake' página 0",
			animeSearch: "naruto hakake",
			page:        "0",
			wantError:   true,
			description: "debe manejar el error al buscar un anime con espacio en el nombre",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scraperClient := animeflv.NewClient()
			ctx := context.Background()

			result, err := scraperClient.SearchAnime(ctx, tc.animeSearch, tc.page)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado en %s: got error = %v, want error = %v", tc.description, err != nil, tc.wantError)
			}

			if !tc.wantError {
				if len(result.Animes) == 0 {
					t.Errorf("esperaba resultados en %s, pero no se encontraron", tc.description)
				}
			}
		})
	}
}

func TestSearchScraper(t *testing.T) {
	testCases := []struct {
		name        string
		wantError   bool
		description string
	}{
		{
			name:        "Prueba búsqueda de todos los animes",
			wantError:   false,
			description: "debe buscar correctamente todos los animes sin errores",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scraperClient := animeflv.NewClient()
			ctx := context.Background()

			_, err := scraperClient.Search(ctx)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado en %s: got error = %v, want error = %v", tc.description, err != nil, tc.wantError)
			}
		})
	}
}

func TestAnimeInfoScraper(t *testing.T) {
	testCases := []struct {
		name        string
		animeID     string
		wantError   bool
		description string
	}{
		{
			name:        "Prueba información de anime con ID 'naruto'",
			animeID:     "naruto",
			wantError:   false,
			description: "debe obtener correctamente la información del anime 'naruto' sin errores",
		},
		{
			name:        "Prueba información de anime con ID 'naruto_shippuden'",
			animeID:     "naruto_shippuden",
			wantError:   true,
			description: "debe manejar el error al obtener información de un anime con ID incorrecto",
		},
		{
			name:        "Prueba información de anime con ID 'one-piece-tv'",
			animeID:     "one-piece-tv",
			wantError:   false,
			description: "debe obtener correctamente la información del anime 'one-piece-tv' sin errores",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scraperClient := animeflv.NewClient()
			ctx := context.Background()

			result, err := scraperClient.AnimeInfo(ctx, tc.animeID)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado en %s: got error = %v, want error = %v", tc.description, err != nil, tc.wantError)
				return
			}

			if !tc.wantError {
				if result.ID != tc.animeID {
					t.Errorf("ID de anime incorrecto: got %s, want %s", result.ID, tc.animeID)
				}
				if len(result.Episodes) == 0 {
					t.Errorf("esperaba episodios en %s, pero no se encontraron", tc.description)
				}
			}
		})
	}
}

func TestLinksScraper(t *testing.T) {
	testCases := []struct {
		name        string
		animeID     string
		episode     uint
		wantError   bool
		description string
	}{
		{
			name:        "Prueba enlaces de episodio 1 del anime 'naruto'",
			animeID:     "naruto",
			episode:     1,
			wantError:   false,
			description: "debe obtener correctamente los enlaces del episodio 1 del anime 'naruto' sin errores",
		},
		{
			name:        "Prueba enlaces de episodio 1154 del anime 'one piece'",
			animeID:     "one-piece-tv",
			episode:     1154,
			wantError:   false,
			description: "debe obtener correctamente los enlaces del episodio 1154 del anime 'one piece' sin errores",
		},
		{
			name:        "Prueba enlaces de episodio 1159 del anime 'one piece'",
			animeID:     "one-piece-tv",
			episode:     1159,
			wantError:   true,
			description: "debe manejar el error al obtener enlaces de un episodio inexistente",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scraperClient := animeflv.NewClient()
			ctx := context.Background()

			_, err := scraperClient.Links(ctx, tc.animeID, tc.episode)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado en %s: got error = %v, want error = %v", tc.description, err != nil, tc.wantError)
			}
		})
	}
}
