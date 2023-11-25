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

	c.wg.Add(1)

	c.scanner.ScanMovies()

	c.wg.Wait()
	lib.Notify(c.ctx.Hub, "import:end", fmt.Sprintf("Import process finished (%s elapsed)", time.Since(t0).String()))
}

func (c *Core) movieFoundListener() {
	for item := range c.movieFoundChan {
		c.wg.Add(1)
		go c.processMovie(item.(*domain.Movie))
	}
}

func (c *Core) processMovie(movie *domain.Movie) {
	defer c.ctx.Hub.Pub(nil, "/event/workunit/done")

	exists := c.storage.CheckExists(movie)
	if exists {
		logger.LightBlue("SKIPPED: exists [%s] (%s)", movie.Title, movie.Location)
		return
	}

	lib.Notify(c.ctx.Hub, "import:progress", fmt.Sprintf("NEW: [%s] (%s)", movie.Title, movie.Location))

	code, err := c.scraper.ScrapeMovie(movie)
	if err != nil {
		logger.Red("SCRAPE MOVIE ERROR: %s", err)
		if code == 2 {
			c.storage.PartialStoreMovie(movie)
		}
		return
	}

	c.storage.StoreMovie(movie)
	c.cache.CacheImages(movie, false)
}

func (c *Core) workUnitDoneListener() {
	for range c.workUnitDone {
		c.wg.Done()
	}
}
