package services

import (
	// sc "bitbucket.org/jdpalmer/statecraft"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	// "io/ioutil"
	"fmt"
	"jbrodriguez/mediagui/server/dto"
	"jbrodriguez/mediagui/server/lib"
	"jbrodriguez/mediagui/server/model"
)

type Core struct {
	Service

	bus      *pubsub.PubSub
	settings *lib.Settings
	// socket   *Socket

	mailbox chan *pubsub.Mailbox

	// m *sc.Machine
}

func NewCore(bus *pubsub.PubSub, settings *lib.Settings) *Core {
	core := &Core{bus: bus, settings: settings}
	core.init()
	return core
}

func (c *Core) Start() {
	mlog.Info("Starting service Core ...")

	c.mailbox = c.register(c.bus, "/get/config", c.getConfig)
	c.registerAdditional(c.bus, "/post/import", c.importMovies, c.mailbox)
	c.registerAdditional(c.bus, "/put/config/folder", c.addMediaFolder, c.mailbox)

	c.registerAdditional(c.bus, "/event/movie/found", c.doMovieFound, c.mailbox)
	c.registerAdditional(c.bus, "/event/movie/scraped", c.doMovieScraped, c.mailbox)

	// c.m = sc.NewMachine("idle")
	// c.m.Rule("import", "idle", "scanning")
	// c.m.Rule("import", "scanning", "scanning")
	// c.m.Rule("found", "scanning", "scanning")
	// c.m.Rule("scraped", "scanning", "scanning")
	// c.m.Rule("status", "idle", "scanning")
	// c.m.Rule("status", "scanning", "scanning")
	// c.m.Rule("finish", "scanning", "idle")

	// data := c.m.Export()
	// ioutil.WriteFile("/Volumes/Users/kayak/tmp/fsm.dot", []byte(data), 0644)

	go c.react()
}

func (c *Core) Stop() {
	mlog.Info("Stopped service Core ...")
}

func (c *Core) react() {
	for mbox := range c.mailbox {
		// mlog.Info("Core:Topic: %s", mbox.Topic)
		c.dispatch(mbox.Topic, mbox.Content)
	}
}

func (c *Core) getConfig(msg *pubsub.Message) {
	msg.Reply <- &c.settings.Config
	mlog.Info("Sent config")
}

func (c *Core) importMovies(msg *pubsub.Message) {
	// mlog.Info("Begin movie scanning ...")

	c.bus.Pub(nil, "/command/movie/scan")
	//	msg.Reply <- &c.settings.Config
	// mlog.Info("Import finished")
}

func (c *Core) addMediaFolder(msg *pubsub.Message) {
	folder := msg.Payload.(string)
	mlog.Info("addMediaFolder: %s", folder)

	c.settings.MediaFolders = append(c.settings.MediaFolders, folder)

	cfg := &pubsub.Message{Payload: c.settings}
	c.bus.Pub(cfg, "/event/config/changed")

	msg.Reply <- &c.settings.Config
	// mlog.Info("Sent config")
}

func (c *Core) doMovieFound(msg *pubsub.Message) {
	movie := msg.Payload.(*model.Movie)

	check := &pubsub.Message{Payload: movie, Reply: make(chan interface{}, 3)}
	c.bus.Pub(check, "/command/movie/exists")

	reply := <-check.Reply
	exists := reply.(bool)

	var text string
	if exists {
		text = fmt.Sprintf("SKIPPED: exists [%s] (%s)", movie.Title, movie.Location)
	} else {
		text = fmt.Sprintf("NEW: [%s] (%s)", movie.Title, movie.Location)
	}

	lib.Notify(c.bus, "import:progress", text)

	if !exists {
		// c.bus.Pub(msg, "/command/movie/scrape")
	}
}

func (c *Core) doMovieScraped(msg *pubsub.Message) {
	dto := msg.Payload.(*dto.Scrape)

	store := &pubsub.Message{Payload: dto.Movie, Reply: make(chan interface{}, 3)}
	c.bus.Pub(store, "/command/movie/store")

	cache := &pubsub.Message{Payload: dto, Reply: make(chan interface{}, 3)}
	c.bus.Pub(cache, "/command/movie/cache")

	mlog.Info("ScrapeDTO: %+v", dto)
}
