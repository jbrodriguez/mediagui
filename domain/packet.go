package domain

// Packet -
type Packet struct {
	ID      string `json:"-"`
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}
