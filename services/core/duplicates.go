package core

import "mediagui/domain"

func (c *Core) GetDuplicates() *domain.MoviesDTO {
	total, items := c.storage.GetDuplicates()
	return &domain.MoviesDTO{Total: total, Items: items}
}
