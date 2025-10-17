package animeflv

import (
	"net/url"
	"strings"

	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/base"
	"github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"
)

type StructAnimeFlv struct {
	base.StructBase
	parser ParserAnimeFlvStruct
}

func AnimeFlvClient() *StructAnimeFlv {
	return &StructAnimeFlv{
		StructBase: base.NewStructBase(
			"https://www3.animeflv.net",
			"https://www3.animeflv.net/browse",
			"https://www3.animeflv.net/anime",
		),
		parser: ParserAnimeFlvStruct{},
	}
}

func ParseURL(baseURL string, params map[string]string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	query := u.Query()
	for key, value := range params {
		query.Add(key, value)
	}

	u.RawQuery = query.Encode()
	return u.String()
}

func (animeflv *StructAnimeFlv) SearchAnime(anime string, page string) ([]dto.SearchAnimeResponse, error) {
	if page == "" {
		page = "1"
	}

	params := map[string]string{
		"page": page,
		"q":    anime,
	}

	url := ParseURL(animeflv.URLSearch(), params)
	resultado, error := animeflv.parser.ParserSearchAnime(url)
	return resultado, error
}

func (animeflv *StructAnimeFlv) AnimeInfo(idAnime string) (dto.ResponseAnimeInfo, error) {
	idAnime = strings.ToLower(idAnime)

	url := animeflv.URLAInfoAnime() + "/" + idAnime
	result, error := animeflv.parser.ParserAnimeInfo(url, idAnime)
	return result, error
}
