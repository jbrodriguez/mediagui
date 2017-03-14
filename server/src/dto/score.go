package dto

type Score struct {
	Id    uint64 `json:"-"`
	Score uint64 `json:"score"`
}
