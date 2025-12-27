// Package main es el punto de entrada de la aplicación API Anime.
// Este archivo contiene ejemplos de uso del servicio de AnimeFlv para demostrar
// todas las funcionalidades disponibles: búsqueda, información de anime, enlaces
// de episodios, animes recientes y episodios recientes.
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	animeflvService "github.com/dst3v3n/api-anime/internal/domain/services/animeflv"
)

func main() {
	// Inicializar el servicio
	service := animeflvService.NewAnimeflvService()

	// Menú de opciones
	if len(os.Args) < 2 {
		printMenu()
		return
	}

	option := os.Args[1]

	switch option {
	case "search":
		testSearchAnime(service)
	case "all-search":
		testSearch(service)
	case "info":
		testAnimeInfo(service)
	case "links":
		testLinks(service)
	case "recent-anime":
		testRecentAnime(service)
	case "recent-episodes":
		testRecentEpisode(service)
	case "all":
		runAllTests(service)
	default:
		fmt.Println("Opción no válida")
		printMenu()
	}
}

// printMenu muestra el menú de opciones disponibles
func printMenu() {
	fmt.Println("=== API Anime - Menu de Pruebas ===")
	fmt.Println("Uso: go run main.go [opcion]")
	fmt.Println("\nOpciones disponibles:")
	fmt.Println("  search           - Buscar anime por nombre con paginación")
	fmt.Println("  all-search       - Obtener todos los animes disponibles")
	fmt.Println("  info             - Información detallada de un anime")
	fmt.Println("  links            - Enlaces de reproducción de un episodio")
	fmt.Println("  recent-anime     - Animes recientemente agregados")
	fmt.Println("  recent-episodes  - Episodios recientemente publicados")
	fmt.Println("  all              - Ejecutar todas las pruebas")
	fmt.Println("\nEjemplo: go run main.go search")
}

// testSearchAnime prueba la búsqueda de anime con paginación
func testSearchAnime(service *animeflvService.AnimeflvService) {
	fmt.Println("\n=== TEST: SearchAnime ===")

	anime := "One Piece"
	page := uint(1)

	fmt.Printf("Buscando: '%s' (Página %d)\n\n", anime, page)

	resultado, err := service.SearchAnime(&anime, &page)
	if err != nil {
		log.Printf("Error en SearchAnime: %v\n", err)
		return
	}

	fmt.Printf("Total de páginas: %d\n", resultado.TotalPages)
	fmt.Printf("Resultados encontrados: %d\n\n", len(resultado.Animes))

	for i, anime := range resultado.Animes {
		fmt.Printf("--- Resultado %d ---\n", i+1)
		fmt.Printf("ID: %s\n", anime.ID)
		fmt.Printf("Título: %s\n", anime.Title)
		fmt.Printf("Tipo: %s\n", anime.Type)
		fmt.Printf("Puntuación: %.1f\n", anime.Punctuation)
		fmt.Printf("Imagen: %s\n", anime.Image)
		fmt.Printf("Sinopsis: %s\n", truncateString(anime.Sinopsis, 100))
		fmt.Println()
	}
}

// testSearch prueba la búsqueda general sin filtros
func testSearch(service *animeflvService.AnimeflvService) {
	fmt.Println("\n=== TEST: Search (Todos los animes) ===\n")

	resultado, err := service.Search()
	if err != nil {
		log.Printf("Error en Search: %v\n", err)
		return
	}

	fmt.Printf("Total de páginas: %d\n", resultado.TotalPages)
	fmt.Printf("Animes encontrados: %d\n\n", len(resultado.Animes))

	// Mostrar solo los primeros 5 resultados para no saturar la consola
	limit := 5
	if len(resultado.Animes) < limit {
		limit = len(resultado.Animes)
	}

	for i := 0; i < limit; i++ {
		anime := resultado.Animes[i]
		fmt.Printf("--- Anime %d ---\n", i+1)
		fmt.Printf("ID: %s\n", anime.ID)
		fmt.Printf("Título: %s\n", anime.Title)
		fmt.Printf("Tipo: %s\n", anime.Type)
		fmt.Printf("Puntuación: %.1f\n", anime.Punctuation)
		fmt.Println()
	}

	if len(resultado.Animes) > limit {
		fmt.Printf("... y %d más\n", len(resultado.Animes)-limit)
	}
}

// testAnimeInfo prueba la obtención de información detallada de un anime
func testAnimeInfo(service *animeflvService.AnimeflvService) {
	fmt.Println("\n=== TEST: AnimeInfo ===")

	idAnime := "one-piece-tv"
	fmt.Printf("Obteniendo información de: %s\n\n", idAnime)

	resultadoInfo, err := service.AnimeInfo(&idAnime)
	if err != nil {
		log.Printf("Error en AnimeInfo: %v\n", err)
		return
	}

	fmt.Println("--- Información General ---")
	fmt.Printf("ID: %s\n", resultadoInfo.ID)
	fmt.Printf("Título: %s\n", resultadoInfo.Title)
	fmt.Printf("Tipo: %s\n", resultadoInfo.Type)
	fmt.Printf("Estado: %s\n", resultadoInfo.Status)
	fmt.Printf("Puntuación: %.1f\n", resultadoInfo.Punctuation)
	fmt.Printf("Imagen: %s\n", resultadoInfo.Image)
	fmt.Printf("Próximo episodio: %s\n\n", resultadoInfo.NextEpisode)

	fmt.Println("--- Géneros ---")
	for _, genre := range resultadoInfo.Genres {
		fmt.Printf("  - %s\n", genre)
	}

	fmt.Println("\n--- Animes Relacionados ---")
	for _, related := range resultadoInfo.AnimeRelated {
		fmt.Printf("  - %s (%s) [ID: %s]\n", related.Title, related.Category, related.ID)
	}

	fmt.Printf("\n--- Episodios ---\n")
	fmt.Printf("Total de episodios: %d\n", len(resultadoInfo.Episodes))
	if len(resultadoInfo.Episodes) > 0 {
		fmt.Printf("Primer episodio: %d\n", resultadoInfo.Episodes[0])
		fmt.Printf("Último episodio: %d\n", resultadoInfo.Episodes[len(resultadoInfo.Episodes)-1])
	}

	fmt.Printf("\n--- Sinopsis ---\n%s\n", resultadoInfo.Sinopsis)
}

// testLinks prueba la obtención de enlaces de reproducción
func testLinks(service *animeflvService.AnimeflvService) {
	fmt.Println("\n=== TEST: Links ===")

	idAnime := "one-piece-tv"
	episode := uint(1146)
	fmt.Printf("Obteniendo enlaces de: %s - Episodio %d\n\n", idAnime, episode)

	resultadoLinks, err := service.Links(&idAnime, &episode)
	if err != nil {
		log.Printf("Error en Links: %v\n", err)
		return
	}

	fmt.Printf("ID: %s\n", resultadoLinks.ID)
	fmt.Printf("Título: %s\n", resultadoLinks.Title)
	fmt.Printf("Episodio: %d\n\n", resultadoLinks.Episode)

	fmt.Println("--- Servidores Disponibles ---")
	for i, link := range resultadoLinks.Link {
		fmt.Printf("\n%d. Servidor: %s\n", i+1, link.Server)
		fmt.Printf("   URL: %s\n", link.URL)
		fmt.Printf("   Code: %s\n", truncateString(link.Code, 50))
	}
}

// testRecentAnime prueba la obtención de animes recientes
func testRecentAnime(service *animeflvService.AnimeflvService) {
	fmt.Println("\n=== TEST: RecentAnime ===\n")

	resultadoRecents, err := service.RecentAnime()
	if err != nil {
		log.Printf("Error en RecentAnime: %v\n", err)
		return
	}

	fmt.Printf("Animes recientes encontrados: %d\n\n", len(resultadoRecents))

	for i, anime := range resultadoRecents {
		fmt.Printf("--- Anime %d ---\n", i+1)
		fmt.Printf("ID: %s\n", anime.ID)
		fmt.Printf("Título: %s\n", anime.Title)
		fmt.Printf("Tipo: %s\n", anime.Type)
		fmt.Printf("Puntuación: %.1f\n", anime.Punctuation)
		fmt.Printf("Imagen: %s\n", anime.Image)
		fmt.Printf("Sinopsis: %s\n", truncateString(anime.Sinopsis, 100))
		fmt.Println()
	}
}

// testRecentEpisode prueba la obtención de episodios recientes
func testRecentEpisode(service *animeflvService.AnimeflvService) {
	fmt.Println("\n=== TEST: RecentEpisode ===\n")

	resultadoEpisodes, err := service.RecentEpisode()
	if err != nil {
		log.Printf("Error en RecentEpisode: %v\n", err)
		return
	}

	fmt.Printf("Episodios recientes encontrados: %d\n\n", len(resultadoEpisodes))

	for i, episode := range resultadoEpisodes {
		fmt.Printf("--- Episodio %d ---\n", i+1)
		fmt.Printf("ID: %s\n", episode.ID)
		fmt.Printf("Título: %s\n", episode.Title)
		fmt.Printf("Capítulo: %s\n", episode.Chapter)
		fmt.Printf("Episodio: %d\n", episode.Episode)
		fmt.Printf("Imagen: %s\n", episode.Image)
		fmt.Println()
	}
}

// runAllTests ejecuta todas las pruebas en secuencia
func runAllTests(service *animeflvService.AnimeflvService) {
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║  Ejecutando todas las pruebas...      ║")
	fmt.Println("╚════════════════════════════════════════╝")

	testSearchAnime(service)
	fmt.Println("\n" + strings.Repeat("=", 50))

	testSearch(service)
	fmt.Println("\n" + strings.Repeat("=", 50))

	testAnimeInfo(service)
	fmt.Println("\n" + strings.Repeat("=", 50))

	testLinks(service)
	fmt.Println("\n" + strings.Repeat("=", 50))

	testRecentAnime(service)
	fmt.Println("\n" + strings.Repeat("=", 50))

	testRecentEpisode(service)

	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║  Todas las pruebas completadas ✓      ║")
	fmt.Println("╚════════════════════════════════════════╝")
}

// truncateString trunca una cadena a una longitud máxima
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
