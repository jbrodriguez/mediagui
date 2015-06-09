package lib

type Options struct {
	Query     string `json:"query"`
	Offset    uint64 `json:"offset"`
	Limit     uint64 `json:"limit"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"`
	FilterBy  string `json:"filterBy"`
}
