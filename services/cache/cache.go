package cache

import (
	"fmt"
	"os"
	"path/filepath"

	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
)

// Cache -
type Cache struct {
	ctx *domain.Context
}

// NewCache -
func Create(ctx *domain.Context) *Cache {
	cache := &Cache{
		ctx: ctx,
	}
	return cache
}

func (c *Cache) Start() {
	logger.Blue("started service cache ...")
}

// Stop -
func (c *Cache) Stop() {
	logger.Blue("stopped service cache")
}

func (c *Cache) CacheImages(movie *domain.Movie, forced bool) {
	defer c.ctx.Hub.Pub(nil, "/event/workunit/done")

	coverPath := filepath.Join(c.ctx.DataDir, "img", "p", movie.Cover)
	if _, err := os.Stat(coverPath); err == nil && !forced {
		lib.Notify(c.ctx.Hub, "import:progress", fmt.Sprintf("COVER DOWNLOAD SKIPPED [%d] %s (%s)", movie.ID, movie.Title, movie.Cover))
	} else {
		if err := lib.Download(movie.SecureBaseURL+"original"+movie.Cover, coverPath); err == nil {
			lib.Notify(c.ctx.Hub, "import:progress", fmt.Sprintf("COVER DOWNLOADED [%d] %s (%s)", movie.ID, movie.Title, movie.Cover))
		} else {
			lib.Notify(c.ctx.Hub, "import:progress", fmt.Sprintf("UNABLE TO DOWNLOAD COVER [%d] %s (%s)", movie.ID, movie.Title, movie.Cover))
		}
	}

	thumbPath := filepath.Join(c.ctx.DataDir, "img", "t", movie.Cover)
	if _, err := os.Stat(thumbPath); err == nil && !forced {
		lib.Notify(c.ctx.Hub, "import:progress", fmt.Sprintf("THUMB GENERATION SKIPPED [%d] %s (%s)", movie.ID, movie.Title, movie.Cover))
	} else {
		if err := lib.ResizeImage(coverPath, thumbPath); err == nil {
			lib.Notify(c.ctx.Hub, "import:progress", fmt.Sprintf("THUMB CREATED [%d] %s (%s)", movie.ID, movie.Title, movie.Cover))
		} else {
			lib.Notify(c.ctx.Hub, "import:progress", fmt.Sprintf("UNABLE TO CREATE THUMB [%d] %s (%s)", movie.ID, movie.Title, movie.Cover))
		}
	}

	backdropPath := filepath.Join(c.ctx.DataDir, "img", "b", movie.Backdrop)
	if _, err := os.Stat(backdropPath); err == nil && !forced {
		lib.Notify(c.ctx.Hub, "import:progress", fmt.Sprintf("BACKDROP DOWNLOAD SKIPPED [%d] %s (%s)", movie.ID, movie.Title, movie.Cover))
	} else {
		if err := lib.Download(movie.SecureBaseURL+"original"+movie.Backdrop, backdropPath); err == nil {
			lib.Notify(c.ctx.Hub, "import:progress", fmt.Sprintf("BACKDROP DOWNLOADED [%d] %s (%s)", movie.ID, movie.Title, movie.Cover))
		} else {
			lib.Notify(c.ctx.Hub, "import:progress", fmt.Sprintf("UNABLE TO DOWNLOAD BACKDROP [%d] %s (%s)", movie.ID, movie.Title, movie.Cover))
		}
	}

}
