package dto

type CategoryAnime string

const (
	Anime    CategoryAnime = "Anime"
	Ova      CategoryAnime = "Ova"
	Pelicula CategoryAnime = "Pelicula"
)

type SearchAnimeResponse struct {
	ID         string
	Title      string
	Sipnopsis  string
	Tipo       CategoryAnime
	Puctuation float64
	Image      string
}
