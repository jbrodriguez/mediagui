package model

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

// ImdbJson -
type ImdbJson struct {
	AggregateRating Rating   `json:"aggregateRating"`
	Director        Entity   `json:"director"`
	Creator         []Entity `json:"creator"`
	Actor           []Entity `json:"actor"`
}

type Imdb struct {
	Director string
	Writers  string
	Actors   string
	Rating   float64
	Votes    uint64
	Awards   string
}
