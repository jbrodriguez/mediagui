package core

import (
	"os"
	"os/signal"
	"syscall"

	"mediagui/domain"
	"mediagui/logger"
	"mediagui/services/cache"
	"mediagui/services/scanner"
	"mediagui/services/scraper"
	"mediagui/services/storage"
)

type Core struct {
	ctx *domain.Context

	storage *storage.Storage
	scanner *scanner.Scanner
	scraper *scraper.Scraper
	cache   *cache.Cache
}

func Create(ctx *domain.Context) *Core {
	return &Core{
		ctx: ctx,

		storage: storage.Create(ctx),
		scanner: scanner.Create(ctx),
		scraper: scraper.Create(ctx),
		cache:   cache.Create(ctx),
	}
}

func (c *Core) Start() error {
	err := c.storage.Start()
	if err != nil {
		return err
	}

	err = c.scraper.Start()
	if err != nil {
		return err
	}

	c.scanner.Start()
	c.cache.Start()

	w := make(chan os.Signal, 1)
	signal.Notify(w, syscall.SIGTERM, syscall.SIGINT)
	logger.Blue("received %s signal. shutting down the app ...", <-w)

	return nil
}
