package core

import "mediagui/domain"

func (c *Core) GetMovies(options *domain.Options) *domain.MoviesDTO {
	total, items := c.storage.GetMovies(options)
	return &domain.MoviesDTO{Total: total, Items: items}
}
