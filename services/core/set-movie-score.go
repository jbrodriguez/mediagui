package core

import (
	"mediagui/domain"
)

func (c *Core) SetMovieScore(dto *domain.MovieDTO) *domain.Movie {
	movie := c.storage.GetMovie(dto.ID)
	movie.Score = dto.Score

	c.storage.SetMovieScore(movie)

	return movie
}
