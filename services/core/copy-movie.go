package core

import "mediagui/domain"

func (c *Core) CopyMovie(dto *domain.FixMovieDTO) *domain.Movie {
	movie := c.storage.GetMovie(dto.ID)
	movie.Tmdb_Id = dto.TmdbID

	return c.storage.CopyMovie(movie)
}
