package main

import (
	"fmt"

	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/animeflv"
)

func main() {
	// Crear una instancia del cliente de AnimeFlv
	client := animeflv.AnimeFlvClient()

	// Buscar un anime específico (por ejemplo "Naruto" en la página 1)
	resultado, _ := client.SearchAnime("One Piece", "1")

	for _, anime := range resultado {
		fmt.Println("ID:", anime.ID)
		fmt.Println("Title:", anime.Title)
		fmt.Println("Sipnopsis:", anime.Sipnopsis)
		fmt.Println("Tipo:", anime.Tipo)
		fmt.Println("Punctuation:", anime.Puctuation)
		fmt.Println("Image:", anime.Image)
		fmt.Println("------------------------------------")
	}

	resultadoInfo, _ := client.AnimeInfo("one-piece-tv")
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
	fmt.Println("Fecha Siguiente Episodio", resultadoInfo.NextEpisode)
}
