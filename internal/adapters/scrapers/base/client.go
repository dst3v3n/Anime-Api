package base

import "github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"

type StructBase struct {
	urlHome       string
	urlSearch     string
	urlAInfoAnime string
}

func (b *StructBase) URLHome() string {
	return b.urlHome
}

func (b *StructBase) URLSearch() string {
	return b.urlSearch
}

func (b *StructBase) URLAInfoAnime() string {
	return b.urlAInfoAnime
}

func NewStructBase(urlHome, urlSearch, urlAInfoAnime string) StructBase {
	return StructBase{
		urlHome:       urlHome,
		urlSearch:     urlSearch,
		urlAInfoAnime: urlAInfoAnime,
	}
}

type MethodBase interface {
	BrowseAnime(anime string)
	GetAnimeInfo(idAnime string)
	GetLink(idAnime string, episode int)
	GetLatestEpisodes()
	GetLatestAnimes()
}

func (b StructBase) SearchAnime(anime string, page int) dto.SearchAnimeResponse {
	return dto.SearchAnimeResponse{}
}

func (b StructBase) AnimeInfo(idAnime string) dto.ResponseAnimeInfo {
	return dto.ResponseAnimeInfo{}
}

func (b StructBase) GetLink(idAnime string, episode int) string {
	return ""
}

func (b StructBase) GetLatestEpisodes() string {
	return ""
}

func (b StructBase) GetLatestAnimes() string {
	return ""
}
