package dto

type Packet struct {
	Id      string `json:"-"`
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}
