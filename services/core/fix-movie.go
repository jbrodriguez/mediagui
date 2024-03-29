package core

import (
	"mediagui/domain"
	"mediagui/logger"
)

func (c *Core) FixMovie(dto *domain.MovieDTO) *domain.Movie {
	movie := c.storage.GetMovie(dto.ID)
	movie.Tmdb_Id = dto.TmdbID

	// 3 operations, rescrape, update and cache
	m, err := c.scraper.ReScrape(movie)
	if err != nil {
		logger.Yellow("RESCRAPE FAILED [%d] %s: %s", movie.ID, movie.Title, err)
		return m
	}

	c.storage.UpdateMovie(m)
	c.cache.CacheImages(movie, true)

	return m
}
