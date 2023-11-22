package storage

import "mediagui/domain"

func (s *Storage) GetMovies(options *domain.Options) (total uint64, items []*domain.Movie) {
	if options.Query == "" {
		total, items = s.listMovies(options)
	} else {
		total, items = s.searchMovies(options)
	}

	return total, items
}
