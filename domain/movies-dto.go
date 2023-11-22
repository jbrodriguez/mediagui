package domain

type MoviesDTO struct {
	Total uint64   `json:"total"`
	Items []*Movie `json:"items"`
}
