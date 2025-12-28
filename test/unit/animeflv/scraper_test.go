// Package animeflv contiene tests unitarios para el scraper de AnimeFlv.
// Este archivo (scraper_test.go) implementa pruebas para verificar el correcto
// funcionamiento del scraper, incluyendo parsing de HTML, extracción de datos
// y manejo de errores. Los tests ayudan a garantizar que los cambios en el
// sitio web sean detectados rápidamente.
package animeflv

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/animeflv"
)

// fixtures HTML embebidos para pruebas //

//go:embed fixtures/search_all.html
var searchAnimeAllHTML []byte

//go:embed fixtures/search_all_fatal.html
var searchAnimeAllFatalHTML []byte

//go:embed fixtures/search_anime.html
var searchAnimeHTML []byte

//go:embed fixtures/search_anime_fatal.html
var searchAnimeFatalHTML []byte

//go:embed fixtures/anime_info.html
var animeInfoHTML []byte

//go:embed fixtures/anime_info_fatal.html
var animeInfoFatalHTML []byte

//go:embed fixtures/episode_anime.html
var episodeLinksHTML []byte

//go:embed fixtures/episode_anime_fatal.html
var episodeLinksFatalHTML []byte

//go:embed fixtures/home_animeflv.html
var homeAnimeflvHTML []byte

//go:embed fixtures/home_animeflv_fatal.html
var homeAnimeflvFatalHTML []byte

func TestParseAnimeSearch(t *testing.T) {
	testCases := []struct {
		name          string
		htmlContent   []byte
		wantError     bool
		expectedCount int
		description   string
	}{
		{
			name:          "búsqueda todos los animes",
			htmlContent:   searchAnimeAllHTML,
			wantError:     false,
			expectedCount: 24,
			description:   "debe parsear correctamente la lista completa de animes",
		},
		{
			name:          "búsqueda todos los animes fatal",
			htmlContent:   searchAnimeAllFatalHTML,
			wantError:     true,
			expectedCount: 0,
			description:   "no deberia parsear correctamente la lista completa de animes debido a error simulado",
		},
		{
			name:          "búsqueda específica de anime",
			htmlContent:   searchAnimeHTML,
			wantError:     false,
			expectedCount: 12,
			description:   "debe parsear correctamente la lista de animes de naruto",
		},
		{
			name:          "búsqueda específica de anime fatal",
			htmlContent:   searchAnimeFatalHTML,
			wantError:     true,
			expectedCount: 0,
			description:   "no deberia parsear correctamente la lista de animes debido a error simulado",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := animeflv.NewParser()
			reader := bytes.NewReader(tc.htmlContent)

			results, err := parser.ParseAnimeWithPagination(reader)

			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado: got %v, want error: %v", err, tc.wantError)
				return
			}

			if !tc.wantError {
				if len(results.Animes) != tc.expectedCount {
					t.Errorf("conteo de resultados incorrecto: got %d, want %d", len(results.Animes), tc.expectedCount)
				}
				for _, anime := range results.Animes {
					if anime.ID == "" || anime.Title == "" {
						t.Errorf("anime con datos incompletos: %+v", anime)
					}
				}
			}
		})
	}
}

func TestParseAnimeInfo(t *testing.T) {
	testCases := []struct {
		name        string
		htmlContent []byte
		wantError   bool
		ID          string
		description string
	}{
		{
			name:        "información de anime exitosa",
			htmlContent: animeInfoHTML,
			wantError:   false,
			ID:          "naruto-shippuden-hd",
			description: "debe parsear correctamente la información del anime Naruto Shippuden",
		},
		{
			name:        "información de anime fatal",
			htmlContent: animeInfoFatalHTML,
			wantError:   true,
			ID:          "naruto-shippuden-hd",
			description: "no deberia parsear correctamente la información del anime debido a error simulado",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := animeflv.NewParser()
			reader := bytes.NewReader(tc.htmlContent)

			result, err := parser.ParseAnimeInfo(reader, tc.ID)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado al parsear información del anime: %v", err)
				return
			}

			if !tc.wantError {
				if result.ID != tc.ID {
					t.Errorf("ID de anime incorrecto: got %s, want %s", result.ID, tc.ID)
				}

				if result.Title == "" || result.Sinopsis == "" {
					t.Errorf("información incompleta del anime: %+v", result)
				}
			}
		})
	}
}

func TestParseLinksEpisode(t *testing.T) {
	testCases := []struct {
		name        string
		htmlContent []byte
		wantError   bool
		description string
	}{
		{
			name:        "enlaces de episodio exitosos",
			htmlContent: episodeLinksHTML,
			wantError:   false,
			description: "debe parsear correctamente los enlaces del episodio",
		},
		{
			name:        "enlaces de episodio fatal",
			htmlContent: episodeLinksFatalHTML,
			wantError:   true,
			description: "no deberia parsear correctamente los enlaces del episodio debido a error simulado",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := animeflv.NewParser()
			reader := bytes.NewReader(tc.htmlContent)

			_, err := parser.ParseLinks(reader, "naruto-shippuden-hd", 220)
			if (err != nil) != tc.wantError {
				t.Errorf("error inesperado al parsear enlaces del episodio: %v", err)
				return
			}
		})
	}
}

func TestParseRecent(t *testing.T) {
	testCases := []struct {
		name        string
		htmlContent []byte
		wantError   bool
	}{
		{
			name:        "episodios recientes exitosos",
			htmlContent: homeAnimeflvHTML,
			wantError:   false,
		},
		{
			name:        "episodios recientes fatal",
			htmlContent: homeAnimeflvFatalHTML,
			wantError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := animeflv.NewParser()

			readerEpisode := bytes.NewReader(tc.htmlContent)
			resultsEpisode, errEpisode := parser.ParseRecentEpisode(readerEpisode)

			if (errEpisode != nil) != tc.wantError {
				t.Errorf("ParseRecentEpisode() error = %v, wantError %v", errEpisode, tc.wantError)
			}

			if !tc.wantError {
				if len(resultsEpisode) == 0 {
					t.Error("ParseRecentEpisode() retornó 0 episodios en caso exitoso")
				} else {
					// Validación adicional del contenido
					firstEp := resultsEpisode[0]
					if firstEp.Title == "" {
						t.Error("ParseRecentEpisode() el primer episodio tiene título vacío")
					}
				}
			}

			readerAnime := bytes.NewReader(tc.htmlContent)
			resultsAnime, errAnime := parser.ParseAnime(readerAnime)

			if (errAnime != nil) != tc.wantError {
				t.Errorf("ParseAnime() error = %v, wantError %v", errAnime, tc.wantError)
			}

			if !tc.wantError {
				if len(resultsAnime) == 0 {
					t.Error("ParseAnime() retornó 0 animes en caso exitoso")
				} else {
					// Validación adicional del contenido
					firstAnime := resultsAnime[0]
					if firstAnime.Title == "" {
						t.Error("ParseAnime() el primer anime tiene título vacío")
					}
				}
			}

			t.Logf("Episodios recientes parseados: %d, Animes recientes parseados: %d",
				len(resultsEpisode), len(resultsAnime))
		})
	}
}
