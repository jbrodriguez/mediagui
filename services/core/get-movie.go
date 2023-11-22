package core

import "mediagui/domain"

func (c *Core) GetMovie(id uint64) *domain.Movie {
	return c.storage.GetMovie(id)
}
