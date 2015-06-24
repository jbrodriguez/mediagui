package lib

type Options struct {
	Query     string `form:"query"`
	Offset    uint64 `form:"offset"`
	Limit     uint64 `form:"limit"`
	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder"`
	FilterBy  string `form:"filterBy"`
}
