package dto

type GetLinkResponse struct {
	ID      string
	Title   string
	Episode int
	Link    []LinkSource
}

type LinkSource struct {
	Server string
	URL    string
	Code   string
}
