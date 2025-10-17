package base

import "github.com/dst3v3n/api-anime/internal/adapters/scrapers/dto"

type BaseParser interface {
	AnimeSearchParser(html string) dto.SearchAnimeResponse
	AnimeInfoParser(html string) dto.ResponseAnimeInfo
}

type ParserFuc struct {
	AnimeSearchParserFunc func(html string) []dto.SearchAnimeResponse
	AnimeInfoParserFunc   func(html string) dto.ResponseAnimeInfo
}

func (p ParserFuc) ParserSearchAnime(url string) ([]dto.SearchAnimeResponse, error) {
	return []dto.SearchAnimeResponse{}, nil
}

func (p ParserFuc) ParserAnimeInfo(url string) (dto.ResponseAnimeInfo, error) {
	return dto.ResponseAnimeInfo{}, nil
}
