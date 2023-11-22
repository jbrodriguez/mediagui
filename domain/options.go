package domain

// Options -
type Options struct {
	Query     string `form:"query" query:"query"`
	Offset    uint64 `form:"offset" query:"offset"`
	Limit     uint64 `form:"limit" query:"limit"`
	SortBy    string `form:"sortBy" query:"sortBy"`
	SortOrder string `form:"sortOrder" query:"sortOrder"`
	FilterBy  string `form:"filterBy" query:"filterBy"`
}
