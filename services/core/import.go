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

	// Reset order coordination for new scan
	c.orderMutex.Lock()
	c.nextSequence = 0
	c.completedMovies = make(map[int]*domain.Movie)
	c.orderMutex.Unlock()

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
			// Mark for partial storage and send to order coordinator
			movie.Stub = 2
			c.orderCoordinatorChan <- movie
		}
		return
	}

	// Send completed movie to order coordinator for sequential DB storage
	c.orderCoordinatorChan <- movie
}

func (c *Core) orderCoordinator() {
	for movie := range c.orderCoordinatorChan {
		c.orderMutex.Lock()

		if movie.Sequence == c.nextSequence {
			// This is the next movie in sequence - store it immediately
			c.storeMovieInOrder(movie)
			c.nextSequence++

			// Check if we can now store any buffered movies
			for {
				if bufferedMovie, exists := c.completedMovies[c.nextSequence]; exists {
					delete(c.completedMovies, c.nextSequence)
					c.storeMovieInOrder(bufferedMovie)
					c.nextSequence++
				} else {
					break
				}
			}
		} else {
			// This movie is out of order - buffer it
			logger.Yellow("Buffering movie seq:%d (waiting for seq:%d): %s", movie.Sequence, c.nextSequence, movie.Title)
			c.completedMovies[movie.Sequence] = movie
		}

		c.orderMutex.Unlock()
	}
}

func (c *Core) storeMovieInOrder(movie *domain.Movie) {
	logger.Green("Storing movie seq:%d in order: %s (%s)", movie.Sequence, movie.Title, movie.Year)

	if movie.Stub == 2 {
		c.storage.PartialStoreMovie(movie)
	} else {
		c.storage.StoreMovie(movie)
		c.cache.CacheImages(movie, false)
	}
}

func (c *Core) workUnitDoneListener() {
	for range c.workUnitDone {
		c.wg.Done()
	}
}
