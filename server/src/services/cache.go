package services

import (
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	// "io/ioutil"
	"fmt"
	// "image/jpeg"
	"jbrodriguez/mediagui/server/src/dto"
	"jbrodriguez/mediagui/server/src/lib"
	"jbrodriguez/mediagui/server/src/model"
	"os"
	"path/filepath"
)

// Cache -
type Cache struct {
	Service

	bus      *pubsub.PubSub
	settings *lib.Settings
	pool     *lib.Pool

	mailbox chan *pubsub.Mailbox
}

// NewCache -
func NewCache(bus *pubsub.PubSub, settings *lib.Settings) *Cache {
	cache := &Cache{bus: bus, settings: settings}
	cache.init()
	return cache
}

// Start -
func (c *Cache) Start() {
	mlog.Info("Starting service Cache ...")

	c.mailbox = c.register(c.bus, "/command/movie/cache", c.cacheMovie)

	c.pool = lib.NewPool(4, 2000)

	go c.react()
}

// Stop -
func (c *Cache) Stop() {
	mlog.Info("Stopped service Cache")
}

func (c *Cache) react() {
	for mbox := range c.mailbox {
		// mlog.Info("Scraper:Topic: %s", mbox.Topic)
		c.dispatch(mbox.Topic, mbox.Content)
	}
}

func (c *Cache) cacheMovie(msg *pubsub.Message) {
	dto := msg.Payload.(*dto.Scrape)
	movie := dto.Movie.(*model.Movie)

	// mlog.Info("trying to cache movie: [%d] %s", movie.Id, movie.Title)

	caching := &Caching{
		bus:      c.bus,
		path:     c.settings.WebDir,
		url:      dto.SecureBaseUrl,
		id:       movie.Id,
		title:    movie.Title,
		forced:   dto.Forced,
		cover:    movie.Cover,
		backdrop: movie.Backdrop,
	}

	c.pool.Exec(caching)
}

// Caching -
type Caching struct {
	bus      *pubsub.PubSub
	path     string
	url      string
	forced   bool
	id       uint64
	title    string
	cover    string
	backdrop string
}

// Execute -
func (c *Caching) Execute(wid int) {
	defer c.bus.Pub(nil, "/event/workunit/done")

	coverPath := filepath.Join(c.path, "img", "p", c.cover)
	if _, err := os.Stat(coverPath); err == nil && !c.forced {
		lib.Notify(c.bus, "import:progress", fmt.Sprintf("COVER DOWNLOAD SKIPPED (%d) [%d] %s (%s)", wid, c.id, c.title, c.cover))
	} else {
		if err := lib.Download(c.url+"original"+c.cover, coverPath); err == nil {
			lib.Notify(c.bus, "import:progress", fmt.Sprintf("COVER DOWNLOADED (%d) [%d] %s (%s)", wid, c.id, c.title, c.cover))
		} else {
			lib.Notify(c.bus, "import:progress", fmt.Sprintf("UNABLE TO DOWNLOAD COVER (%d) [%d] %s (%s)", wid, c.id, c.title, c.cover))
		}
	}

	thumbPath := filepath.Join(c.path, "img", "t", c.cover)
	if _, err := os.Stat(thumbPath); err == nil && !c.forced {
		lib.Notify(c.bus, "import:progress", fmt.Sprintf("THUMB GENERATION SKIPPED (%d) [%d] %s (%s)", wid, c.id, c.title, c.cover))
	} else {
		if err := lib.ResizeImage(coverPath, thumbPath); err == nil {
			lib.Notify(c.bus, "import:progress", fmt.Sprintf("THUMB CREATED (%d) [%d] %s (%s)", wid, c.id, c.title, c.cover))
		} else {
			lib.Notify(c.bus, "import:progress", fmt.Sprintf("UNABLE TO CREATE THUMB (%d) [%d] %s (%s)", wid, c.id, c.title, c.cover))
		}
	}

	backdropPath := filepath.Join(c.path, "img", "b", c.backdrop)
	if _, err := os.Stat(backdropPath); err == nil && !c.forced {
		lib.Notify(c.bus, "import:progress", fmt.Sprintf("BACKDROP DOWNLOAD SKIPPED (%d) [%d] %s (%s)", wid, c.id, c.title, c.cover))
	} else {
		if err := lib.Download(c.url+"original"+c.backdrop, backdropPath); err == nil {
			lib.Notify(c.bus, "import:progress", fmt.Sprintf("BACKDROP DOWNLOADED (%d) [%d] %s (%s)", wid, c.id, c.title, c.cover))
		} else {
			lib.Notify(c.bus, "import:progress", fmt.Sprintf("UNABLE TO DOWNLOAD BACKDROP (%d) [%d] %s (%s)", wid, c.id, c.title, c.cover))
		}
	}

	// event := "/event/movie/cached"
	// if c.forced {
	// 	event += "/forced"
	// }

	// cached := &pubsub.Message{}
	// c.bus.Pub(cached, event)
}