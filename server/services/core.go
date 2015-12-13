package services

import (
	// sc "bitbucket.org/jdpalmer/statecraft"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/myodc/go-micro/client"
	// "io/ioutil"
	"fmt"
	"jbrodriguez/mediagui/proto"
	"jbrodriguez/mediagui/server/dto"
	"jbrodriguez/mediagui/server/lib"
	"jbrodriguez/mediagui/server/model"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
)

type Core struct {
	Service

	bus      *pubsub.PubSub
	settings *lib.Settings
	// socket   *Socket

	mailbox chan *pubsub.Mailbox
	re      *regexp.Regexp
	maps    map[string]bool

	wg sync.WaitGroup
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
	c.registerAdditional(c.bus, "/post/prune", c.pruneMovies, c.mailbox)
	c.registerAdditional(c.bus, "/put/config/folder", c.addMediaFolder, c.mailbox)
	c.registerAdditional(c.bus, "/put/movies/fix", c.fixMovie, c.mailbox)

	c.registerAdditional(c.bus, "/event/movie/found", c.doMovieFound, c.mailbox)
	c.registerAdditional(c.bus, "/event/movie/tmdbnotfound", c.doMovieTmdbNotFound, c.mailbox)
	c.registerAdditional(c.bus, "/event/movie/scraped", c.doMovieScraped, c.mailbox)
	// c.registerAdditional(c.bus, "/event/movie/rescraped", c.doMovieReScraped, c.mailbox)
	// c.registerAdditional(c.bus, "/event/movie/updated", c.doMovieUpdated, c.mailbox)
	// c.registerAdditional(c.bus, "/event/movie/cached/forced", c.doMovieCachedForced, c.mailbox)

	c.registerAdditional(c.bus, "/event/workunit/done", c.doWorkUnitDone, c.mailbox)

	c.re = regexp.MustCompile(`(?i)/Volumes/(.*?)/.*`)
	c.maps = make(map[string]bool)

	// maps["/Volumes/wopr-films"] = true // for example
	for _, folder := range c.settings.MediaFolders {
		c.maps[folder] = false
		if _, err := os.Stat(folder); err == nil {
			c.maps[folder] = true
		}
	}

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
	t0 := time.Now()
	// mlog.Info("Begin movie scanning ...")
	lib.Notify(c.bus, "import:begin", "Import process started")

	c.wg.Add(1)

	c.bus.Pub(nil, "/command/movie/scan")
	//	msg.Reply <- &c.settings.Config
	// mlog.Info("Import finished")

	go func() {
		c.wg.Wait()
		// 𝛥t := float64(time.Since(t0)) / 1e9
		lib.Notify(c.bus, "import:end", fmt.Sprintf("Import process finished (%s elapsed)", time.Since(t0).String()))
	}()

}

func (c *Core) doMovieFound(msg *pubsub.Message) {
	movie := msg.Payload.(*model.Movie)

	check := &pubsub.Message{Payload: movie, Reply: make(chan interface{}, 3)}
	c.bus.Pub(check, "/command/movie/exists")

	reply := <-check.Reply
	exists := reply.(bool)

	if exists {
		mlog.Info("SKIPPED: exists [%s] (%s)", movie.Title, movie.Location)
	} else {
		lib.Notify(c.bus, "import:progress", fmt.Sprintf("NEW: [%s] (%s)", movie.Title, movie.Location))

		c.wg.Add(1)
		c.bus.Pub(msg, "/command/movie/scrape")
	}
}

func (c *Core) fixMovie(msg *pubsub.Message) {
	movie := msg.Payload.(*model.Movie)

	// 3 operations, rescrape, update and cache
	c.wg.Add(1)

	// rescrape
	scrape := &pubsub.Message{Payload: movie, Reply: make(chan interface{}, 3)}
	c.bus.Pub(scrape, "/command/movie/rescrape")

	go func() {
		c.wg.Wait()
		msg.Reply <- movie
	}()
	// go c.waitFixMovie(msg.Reply, movie)
}

func (c *Core) doMovieTmdbNotFound(msg *pubsub.Message) {
	dto := msg.Payload.(*dto.Scrape)

	store := &pubsub.Message{Payload: dto.Movie, Reply: make(chan interface{}, 3)}
	c.bus.Pub(store, "/command/movie/partialstore")
}

func (c *Core) doMovieScraped(msg *pubsub.Message) {
	dto := msg.Payload.(*dto.Scrape)

	mlog.Info("ScrapeDTO: %+v", dto)

	// I treat the following two commands as one, for the sake of the wg
	// now there are two oustanding locks, which will be decreased by each
	// responding service
	c.wg.Add(1)

	store := &pubsub.Message{Payload: dto.Movie, Reply: make(chan interface{}, 3)}
	c.bus.Pub(store, "/command/movie/store")

	cache := &pubsub.Message{Payload: dto, Reply: make(chan interface{}, 3)}
	c.bus.Pub(cache, "/command/movie/cache")

	// mlog.Info("ScrapeDTO: %+v", dto)
}

func (c *Core) pruneMovies(msg *pubsub.Message) {
	t0 := time.Now()

	lib.Notify(c.bus, "prune:begin", "Started Prune Process")

	options := &lib.Options{Offset: 0, Limit: 99999999999999, SortBy: "title", SortOrder: "asc"}
	all := &pubsub.Message{Payload: options, Reply: make(chan interface{}, capacity)}
	c.bus.Pub(all, "/get/movies")

	reply := <-all.Reply
	dto := reply.(*model.MoviesDTO)

	if c.settings.UnraidMode {
		for _, item := range dto.Items {
			// mlog.Info("Item is %s (%s)", item.Title, item.Location)

			index := strings.Index(item.Location, ":")

			host := item.Location[:index]
			location := item.Location[index+1:]

			req := client.NewRequest("io.jbrodriguez.mediagui.agent."+host, "Agent.Exists", &agent.ExistsReq{
				Location: location,
			})

			rsp := &agent.ExistsRsp{}

			if err := client.Call(context.Background(), req, rsp); err != nil {
				mlog.Warning("Unable to connect to service (%s): %s", "io.jbrodriguez.mediagui.agent."+host, err)
				// lib.Notify(s.bus, "import:progress", "Unable to connect to host "+host)
				// lib.Notify(s.bus, "import:end", "Import process finished")
				// return
				continue
			}

			if !rsp.Exists {
				lib.Notify(c.bus, "prune:selected", fmt.Sprintf("UP FOR DELETION: [%d] %s (%s))", item.Id, item.Title, item.Location))

				movie := &pubsub.Message{Payload: item, Reply: make(chan interface{}, capacity)}
				c.bus.Pub(movie, "/command/movie/delete")
			}
		}
	} else {
		for _, item := range dto.Items {
			matches := c.re.FindStringSubmatch(item.Location)
			if len(matches) == 0 {
				continue
			}

			folder := filepath.Join("/Volumes", matches[1])
			if !c.maps[folder] {
				// mlog.Info("Folder not mapped (%s): %s", folder, item.Location)
				continue
			}

			// mlog.Info("Folder mapped (%s): %s", folder, item.Location)

			if _, err := os.Stat(item.Location); err != nil {
				if os.IsNotExist(err) {
					lib.Notify(c.bus, "prune:selected", fmt.Sprintf("UP FOR DELETION: [%d] %s (%s))", item.Id, item.Title, item.Location))

					movie := &pubsub.Message{Payload: item, Reply: make(chan interface{}, capacity)}
					c.bus.Pub(movie, "/command/movie/delete")
				}
			}
		}
	}

	// for _, folder := range c.settings.MediaFolders {
	// 	mlog.Info("folder: %s", folder)
	// 	if _, err := os.Stat(folder); err != nil {
	// 		if os.IsNotExist(err) {
	// 			lib.Notify(c.bus, "prune:error", fmt.Sprintf("Folder %s is not present. Aborting Prune process.", folder))
	// 			return
	// 		}
	// 	}
	// }

	// lib.Notify(c.bus, "prune:end", "Finished Prune Process")
	lib.Notify(c.bus, "prune:end", fmt.Sprintf("Prune process finished (%s elapsed)", time.Since(t0).String()))

}

func (c *Core) addMediaFolder(msg *pubsub.Message) {
	folder := msg.Payload.(string)
	mlog.Info("addMediaFolder: %s", folder)

	c.settings.MediaFolders = append(c.settings.MediaFolders, folder)
	c.settings.Save()

	cfg := &pubsub.Message{Payload: c.settings}
	c.bus.Pub(cfg, "/event/config/changed")

	msg.Reply <- &c.settings.Config
	// mlog.Info("Sent config")
}

// func (c *Core) waitFixMovie(ch chan interface{}, movie *model.Movie) {
// 	c.wg.Wait()
// 	ch <- movie
// }

// func (c *Core) doMovieUpdated(msg *pubsub.Message) {
// 	c.wg.Done()
// }

// func (c *Core) doMovieCachedForced(msg *pubsub.Message) {
// 	c.wg.Done()
// }

func (c *Core) doWorkUnitDone(msg *pubsub.Message) {
	c.wg.Done()
}
