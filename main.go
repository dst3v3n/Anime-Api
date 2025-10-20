// Package main es el punto de entrada de la aplicación API Anime.
// Este archivo contiene ejemplos de uso del scraper de AnimeFlv para demostrar
// todas las funcionalidades disponibles: búsqueda, información de anime, enlaces
// de episodios, animes recientes y episodios recientes.
package main

import (
	"fmt"

	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/animeflv"
	"github.com/dst3v3n/api-anime/internal/ports"
)

func main() {
	var scraper ports.ScraperPort = animeflv.NewClient()

	resultado, err := scraper.SearchAnime("One Piece", "1")
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("=================Resultado=================")
	for _, anime := range resultado {
		fmt.Println("ID:", anime.ID)
		fmt.Println("Titulo:", anime.Title)
		fmt.Println("Imagen:", anime.Image)
		fmt.Println("Descripcion:", anime.Sinopsis)
		fmt.Println("Tipo", anime.Type)
		fmt.Println("Punctuacion:", anime.Punctuation)
		fmt.Println("------------------------------------------")
	}

	resultadoInfo, err := scraper.AnimeInfo("one-piece-tv")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("===============Anime Info=================")
	fmt.Println("ID:", resultadoInfo.ID)
	fmt.Println("Titulo:", resultadoInfo.Title)
	fmt.Println("Imagen:", resultadoInfo.Image)
	fmt.Println("Descripcion:", resultadoInfo.Sinopsis)
	fmt.Println("Tipo", resultadoInfo.Type)
	fmt.Println("Punctuacion:", resultadoInfo.Punctuation)
	fmt.Println("Estado:", resultadoInfo.Status)
	fmt.Println("Generos:", resultadoInfo.Genres)
	fmt.Println("Episodios:", resultadoInfo.Episodes)
	fmt.Println("Proximo Episodio:", resultadoInfo.NextEpisode)

	resultadoLinks, err := scraper.Links("one-piece-tv", 1150)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("===============Links Episodio=================")
	fmt.Println("ID:", resultadoLinks.ID)
	fmt.Println("Titulo:", resultadoLinks.Title)
	fmt.Println("Episodio:", resultadoLinks.Episode)
	for _, link := range resultadoLinks.Link {
		fmt.Println("Servidor:", link.Server)
		fmt.Println("URL:", link.URL)
		fmt.Println("Code:", link.Code)
	}

	resultadoRecents, err := scraper.RecentAnime()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("===============Animes Recientes=================")
	for _, anime := range resultadoRecents {
		fmt.Println("ID:", anime.ID)
		fmt.Println("Titulo:", anime.Title)
		fmt.Println("Imagen:", anime.Image)
		fmt.Println("Descripcion:", anime.Sinopsis)
		fmt.Println("Tipo", anime.Type)
		fmt.Println("Punctuacion:", anime.Punctuation)
		fmt.Println("------------------------------------------")
	}

	resultadoEpisodes, err := scraper.RecentEpisode()
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("===============Episodios Recientes=================")
	for _, episode := range resultadoEpisodes {
		fmt.Println("ID:", episode.ID)
		fmt.Println("Titulo:", episode.Title)
		fmt.Println("Capitulo:", episode.Chapter)
		fmt.Println("Episodio:", episode.Episode)
		fmt.Println("Imagen:", episode.Image)
		fmt.Println("------------------------------------------")
	}
}
