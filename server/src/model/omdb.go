package model

// Omdb -
type Omdb struct {
	Director    string `json:"Director"`
	Writer      string `json:"Writer"`
	Actors      string `json:"Actors"`
	Awards      string `json:"Awards"`
	Imdb_Rating string `json:"imdbRating"`
	Imdb_Vote   string `json:"imdbVotes"`
}
