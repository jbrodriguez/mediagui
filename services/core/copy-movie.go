package core

import "mediagui/domain"

func (c *Core) CopyMovie(movie *domain.Movie) *domain.Movie {
	return c.storage.CopyMovie(movie)
}
