package core

import "mediagui/domain"

func (c *Core) SetDuplicate(movie *domain.Movie) *domain.Movie {
	m := c.storage.SetDuplicate(movie)
	return m
}
