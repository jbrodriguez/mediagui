package domain

type MovieDTO struct {
	ID      uint64 `json:"id"`
	TmdbID  uint64 `json:"tmdb_id"`
	Score   uint64 `json:"score"`
	Watched string `json:"watched"`
}
