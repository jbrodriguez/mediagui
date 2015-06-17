package services

import (
	// sc "bitbucket.org/jdpalmer/statecraft"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	// "io/ioutil"
	"fmt"
	"jbrodriguez/mediagui/server/lib"
	"jbrodriguez/mediagui/server/model"
)

type Scraper struct {
	Service

	bus      *pubsub.PubSub
	settings *lib.Settings
	// socket   *Socket

	mailbox chan *pubsub.Mailbox

	// m *sc.Machine
}

func NewScraper(bus *pubsub.PubSub, settings *lib.Settings) *Scraper {
	scraper := &Scraper{bus: bus, settings: settings}
	scraper.init()
	return scraper
}

func (s *Scraper) Start() {
	mlog.Info("Starting service Scraper ...")

	s.mailbox = s.register(s.bus, "/command/scrape/movie", s.scrapeMovie)

	go s.react()
}

func (s *Scraper) Stop() {
	mlog.Info("Stopped service Scraper ...")
}

func (s *Scraper) react() {
	for mbox := range s.mailbox {
		// mlog.Info("Scraper:Topic: %s", mbox.Topic)
		s.dispatch(mbox.Topic, mbox.Content)
	}
}
