package storage

import (
	"mediagui/domain"
)

func (s *Storage) GetCovers() (total uint64, items []*domain.Movie) {
	options := &domain.Options{
		Limit:     60,
		SortBy:    "added",
		SortOrder: "desc",
	}

	return s.listMovies(options)
}
