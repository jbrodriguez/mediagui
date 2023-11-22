package domain

import (
	"encoding/json"
)

// Entity -
type EntityAlt struct {
	Type string `json:"@type"`
	URL  string `json:"url"`
	Name string `json:"name"`
}

// Rating -
type RatingAlt struct {
	Type        string  `json:"@type"`
	RatingCount int     `json:"ratingCount"`
	RatingValue float64 `json:"ratingValue"`
}

// ImdbJSON -
type ImdbJSONAlt struct {
	AggregateRating RatingAlt       `json:"aggregateRating"`
	RawDirector     json.RawMessage `json:"director"`
	Director        []EntityAlt     `json:"-"`
	RawCreator      json.RawMessage `json:"creator"`
	Creator         []EntityAlt     `json:"-"`
	Actor           []EntityAlt     `json:"actor"`
}

const PERSONAlt = "Person"

// https://blog.gopheracademy.com/advent-2016/advanced-encoding-decoding/

// Imdb - Creates an Imdb object from its go/json form
func (ij *ImdbJSONAlt) ImdbAlt() ImdbAlt {
	imdb := ImdbAlt{
		Votes:  uint64(ij.AggregateRating.RatingCount),
		Rating: ij.AggregateRating.RatingValue,
	}

	for _, director := range ij.Director {
		if director.Type == PERSONAlt {
			if imdb.Director == "" {
				imdb.Director = director.Name
			} else {
				imdb.Director += ", " + director.Name
			}
		}
	}

	for _, writer := range ij.Creator {
		if writer.Type == PERSONAlt {
			if imdb.Writers == "" {
				imdb.Writers = writer.Name
			} else {
				imdb.Writers += ", " + writer.Name
			}
		}
	}

	for _, actor := range ij.Actor {
		if actor.Type == PERSONAlt {
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
type ImdbAlt struct {
	Director string
	Writers  string
	Actors   string
	Rating   float64
	Votes    uint64
	Awards   string
}

// UnmarshalJSON -
func (i *ImdbAlt) UnmarshalJSON(data []byte) error {
	var ij ImdbJSONAlt

	if err := json.Unmarshal(data, &ij); err != nil {
		return err
	}

	if len(ij.RawCreator) > 0 {
		switch ij.RawCreator[0] {
		case '{':
			var entity EntityAlt
			if err := json.Unmarshal(ij.RawCreator, &entity); err != nil {
				return err
			}
			ij.Creator = []EntityAlt{entity}
		case '[':
			if err := json.Unmarshal(ij.RawCreator, &ij.Creator); err != nil {
				return err
			}
		}
	}

	if len(ij.RawDirector) > 0 {
		switch ij.RawDirector[0] {
		case '{':
			var entity EntityAlt
			if err := json.Unmarshal(ij.RawDirector, &entity); err != nil {
				return err
			}
			ij.Director = []EntityAlt{entity}
		case '[':
			if err := json.Unmarshal(ij.RawDirector, &ij.Director); err != nil {
				return err
			}
		}
	}

	*i = ij.ImdbAlt()

	return nil
}
