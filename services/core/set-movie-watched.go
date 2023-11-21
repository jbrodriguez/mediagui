package core

import "mediagui/domain"

func (c *Core) SetMovieWatched(dto *domain.MovieDTO) *domain.Movie {
	movie := c.storage.GetMovie(dto.ID)
	movie.Last_Watched = dto.Watched

	return c.storage.SetMovieWatched(movie)
}
