package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type OmdbResponse struct {
	Title      string       `json:"Title"`
	Year       string       `json:"Year"`
	Rated      string       `json:"Rated"`
	Released   string       `json:"Released"`
	Runtime    string       `json:"Runtime"`
	Genre      string       `json:"Genre"`
	Director   string       `json:"Director"`
	Writer     string       `json:"Writer"`
	Actors     string       `json:"Actors"`
	Plot       string       `json:"Plot"`
	Language   string       `json:"Language"`
	Country    string       `json:"Country"`
	Awards     string       `json:"Awards"`
	Poster     string       `json:"Poster"`
	Ratings    []OmdbRating `json:"Ratings"`
	Metascore  string       `json:"Metascore"`
	ImdbRating string       `json:"imdbRating"`
	ImdbVotes  string       `json:"imdbVotes"`
	ImdbID     string       `json:"imdbID"`
	Type       string       `json:"Type"`
	Response   string       `json:"Response"`
	Error      string       `json:"Error"`
}

type OmdbRating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

func (o *OmdbResponse) GetImdbRating() float64 {
	rating, _ := strconv.ParseFloat(o.ImdbRating, 64)
	return rating
}

func (o *OmdbResponse) GetImdbVotes() uint64 {
	cleaned := strings.ReplaceAll(o.ImdbVotes, ",", "")
	votes, _ := strconv.ParseUint(cleaned, 10, 64)
	return votes
}

func FetchOmdb(imdbID string, apiKey string) (*OmdbResponse, error) {
	url := fmt.Sprintf("http://www.omdbapi.com/?i=%s&apikey=%s", imdbID, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OMDb data: %w", err)
	}
	defer resp.Body.Close()

	var omdb OmdbResponse
	if err := json.NewDecoder(resp.Body).Decode(&omdb); err != nil {
		return nil, fmt.Errorf("failed to decode OMDb response: %w", err)
	}

	if omdb.Response == "False" {
		return nil, fmt.Errorf("OMDb error: %s", omdb.Error)
	}

	return &omdb, nil
}
