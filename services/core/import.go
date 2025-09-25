package core

import (
	"fmt"
	"time"

	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
)

func (c *Core) ImportMovies() {
	t0 := time.Now()
	lib.Notify(c.ctx.Hub, "import:begin", "Import process started")

	// Get all movies from scanner in scan order
	movies := c.scanner.ScanMovies()

	// Process each movie sequentially
	for i, movie := range movies {
		logger.Blue("Processing movie %d/%d: %s (%s)", i+1, len(movies), movie.Title, movie.Year)
		c.processMovie(movie)
	}

	lib.Notify(c.ctx.Hub, "import:end", fmt.Sprintf("Import process finished (%s elapsed)", time.Since(t0).String()))
}

func (c *Core) processMovie(movie *domain.Movie) bool {
	exists := c.storage.CheckExists(movie)
	if exists {
		logger.LightBlue("SKIPPED: exists [%s] (%s)", movie.Title, movie.Location)
		return true
	}

	lib.Notify(c.ctx.Hub, "import:progress", fmt.Sprintf("NEW: [%s] (%s)", movie.Title, movie.Location))

	code, err := c.scraper.ScrapeMovie(movie)
	if err != nil {
		logger.Red("SCRAPE MOVIE ERROR: %s", err)
		if code == 2 {
			c.storage.PartialStoreMovie(movie)
		}
		return false
	}

	c.storage.StoreMovie(movie)
	c.cache.CacheImages(movie, false)

	return false
}
