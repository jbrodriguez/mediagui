package core

import "mediagui/domain"

func (c *Core) GetMovie(id string) *domain.Movie {
	return c.storage.GetMovie(id)
}
