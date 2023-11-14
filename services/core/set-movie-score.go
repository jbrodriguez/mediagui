package core

import (
	"mediagui/domain"
)

func (c *Core) SetMovieScore(movie *domain.Movie) *domain.Movie {
	c.storage.SetMovieScore(movie)
	return movie
}
