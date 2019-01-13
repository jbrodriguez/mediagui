package model

import (
	"encoding/json"
	"strconv"
)

// Entity -
type Entity struct {
	Type string `json:"@type"`
	URL  string `json:"url"`
	Name string `json:"name"`
}

// Rating -
type Rating struct {
	Type        string `json:"@type"`
	RatingCount int    `json:"ratingCount"`
	RatingValue string `json:"ratingValue"`
}

// ImdbJSON -
type ImdbJSON struct {
	AggregateRating Rating          `json:"aggregateRating"`
	RawDirector     json.RawMessage `json:"director"`
	Director        []Entity        `json:"-"`
	RawCreator      json.RawMessage `json:"creator"`
	Creator         []Entity        `json:"-"`
	Actor           []Entity        `json:"actor"`
}

// https://blog.gopheracademy.com/advent-2016/advanced-encoding-decoding/

// Imdb - Creates an Imdb object from its go/json form
func (ij ImdbJSON) Imdb() Imdb {
	rating, _ := strconv.ParseFloat(ij.AggregateRating.RatingValue, 64)

	imdb := Imdb{
		Votes:    uint64(ij.AggregateRating.RatingCount),
		Rating:   rating,
	}

	for _, director := range ij.Director {
		if director.Type == "Person" {
			if imdb.Director == "" {
				imdb.Director = director.Name
			} else {
				imdb.Director += ", " + director.Name
			}
		}
	}

	for _, writer := range ij.Creator {
		if writer.Type == "Person" {
			if imdb.Writers == "" {
				imdb.Writers = writer.Name
			} else {
				imdb.Writers += ", " + writer.Name
			}
		}
	}

	for _, actor := range ij.Actor {
		if actor.Type == "Person" {
			if imdb.Actors == "" {
				imdb.Actors = actor.Name
			} else {
				imdb.Actors += ", " + actor.Name
			}
		}
	}

	return imdb
}

// Imdb -
type Imdb struct {
	Director string
	Writers  string
	Actors   string
	Rating   float64
	Votes    uint64
	Awards   string
}

// UnmarshalJSON -
func (i *Imdb) UnmarshalJSON(data []byte) error {
	var ij ImdbJSON

	if err := json.Unmarshal(data, &ij); err != nil {
		return err
	}

	if len(ij.RawCreator) > 0 {
		switch ij.RawCreator[0] {
		case '{':
			var entity Entity
			if err := json.Unmarshal(ij.RawCreator, &entity); err != nil {
				return err
			}
			ij.Creator = []Entity{entity}
		case '[':
			if err := json.Unmarshal(ij.RawCreator, &ij.Creator); err != nil {
				return err
			}
		}
	}

	if len(ij.RawDirector) > 0 {
		switch ij.RawDirector[0] {
		case '{':
			var entity Entity
			if err := json.Unmarshal(ij.RawDirector, &entity); err != nil {
				return err
			}
			ij.Director = []Entity{entity}
		case '[':
			if err := json.Unmarshal(ij.RawDirector, &ij.Director); err != nil {
				return err
			}
		}
	}

	*i = ij.Imdb()

	return nil
}
