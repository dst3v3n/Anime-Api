package dto

type StatusAnime string

const (
	Emision    StatusAnime = "En Emision"
	Finalizado StatusAnime = "Finalizado"
)

type ResponseAnimeInfo struct {
	SearchAnimeResponse
	AnimeRelated []AnimeRelated
	Generos      []string
	Estado       StatusAnime
	NextEpisode  string
	Episodes     []int
}

type AnimeRelated struct {
	ID       string
	Title    string
	Category string
}
