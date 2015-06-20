package model

type ScrapeDTO struct {
	BaseUrl       string
	SecureBaseUrl string
	BasePath      string
	Movie         *Movie
	Forced        bool
}

type Packet struct {
	Id      string `json:"-"`
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}
