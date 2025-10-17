package main

import (
	"fmt"

	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/animeflv"
	"github.com/dst3v3n/api-anime/internal/ports"
)

func main() {
	// Usar la interfaz ScraperPort en lugar de la implementación concreta
	var scraper ports.ScraperPort = animeflv.NewClient()

	// Buscar un anime específico (por ejemplo "One Piece" en la página 1)
	resultado, err := scraper.SearchAnime("One Piece", "1")
	if err != nil {
		fmt.Println("Error al buscar anime:", err)
		return
	}

	fmt.Println("=== RESULTADOS DE BÚSQUEDA ===")
	for _, anime := range resultado {
		fmt.Println("ID:", anime.ID)
		fmt.Println("Title:", anime.Title)
		fmt.Println("Sipnopsis:", anime.Sipnopsis)
		fmt.Println("Tipo:", anime.Tipo)
		fmt.Println("Punctuation:", anime.Puctuation)
		fmt.Println("Image:", anime.Image)
		fmt.Println("------------------------------------")
	}

	// // Obtener información detallada de un anime
	resultadoInfo, err := scraper.AnimeInfo("one-piece-tv")
	if err != nil {
		fmt.Println("Error al obtener info del anime:", err)
		return
	}

	fmt.Println("\n=== INFORMACIÓN DETALLADA ===")
	fmt.Println("ID:", resultadoInfo.ID)
	fmt.Println("Title:", resultadoInfo.Title)
	fmt.Println("Sipnopsis:", resultadoInfo.Sipnopsis)
	fmt.Println("Tipo:", resultadoInfo.Tipo)
	fmt.Println("Puntuacion:", resultadoInfo.Puctuation)
	fmt.Println("Image:", resultadoInfo.Image)
	fmt.Println("Animes Relacionados:", resultadoInfo.AnimeRelated)
	fmt.Println("Genres:", resultadoInfo.Generos)
	fmt.Println("Status:", resultadoInfo.Estado)
	fmt.Println("Episodes:", resultadoInfo.Episodes)
	fmt.Println("Fecha Siguiente Episodio:", resultadoInfo.NextEpisode)

	linksEpidosde, errorEpisode := scraper.GetLinks("one-piece-tv", 1145)
	if errorEpisode != nil {
		fmt.Println("Error al obtener links del episodio:", errorEpisode)
		return
	}
	fmt.Println("=== LINKS DEL EPISODIO ===")
	fmt.Printf("ID : %s \n", linksEpidosde.ID)
	fmt.Printf("Título del episodio: %s \n", linksEpidosde.Title)
	fmt.Printf("Episodio: %d \n", linksEpidosde.Episode)
	fmt.Printf("Links del episodio: %+v \n", linksEpidosde)
}
