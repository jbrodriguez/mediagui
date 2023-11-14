package core

import "mediagui/domain"

func (c *Core) SetMovieWatched(movie *domain.Movie) *domain.Movie {
	return c.storage.SetMovieWatched(movie)
}
