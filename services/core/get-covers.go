package core

import "mediagui/domain"

func (c *Core) GetCovers() *domain.MoviesDTO {
	total, items := c.storage.GetCovers()
	return &domain.MoviesDTO{Total: total, Items: items}
}
