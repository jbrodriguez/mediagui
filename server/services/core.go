package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"google.golang.org/grpc"

	"mediagui/dto"
	"mediagui/lib"
	pb "mediagui/mediaagent"
	"mediagui/model"
)

const cNotAvailable = "n/a"

// Core -
type Core struct {
	bus      *pubsub.PubSub
	settings *lib.Settings

	actor *actor.Actor
	re    *regexp.Regexp
	maps  map[string]bool

	wg sync.WaitGroup
}

// NewCore -
func NewCore(bus *pubsub.PubSub, settings *lib.Settings) *Core {
	core := &Core{
		bus:      bus,
		settings: settings,
		actor:    actor.NewActor(bus),
	}
	return core
}

// Start -
func (c *Core) Start() {
	mlog.Info("Starting service Core ...")

	c.actor.Register("/get/config", c.getConfig)
	c.actor.Register("/post/import", c.importMovies)
	c.actor.Register("/post/add", c.addMovie)
	c.actor.Register("/post/prune", c.pruneMovies)
	c.actor.Register("/put/config/folder", c.addMediaFolder)
	c.actor.Register("/put/movies/fix", c.fixMovie)

	c.actor.Register("/event/movie/found", c.doMovieFound)
	c.actor.Register("/event/movie/tmdbnotfound", c.doMovieTmdbNotFound)
	c.actor.Register("/event/movie/scraped", c.doMovieScraped)
	c.actor.Register("/event/movie/rescraped", c.doMovieReScraped)

	c.actor.Register("/event/workunit/done", c.doWorkUnitDone)

	c.re = regexp.MustCompile(`(?i)/Volumes/(.*?)/.*`)
	c.maps = make(map[string]bool)

	// for example, maps["/Volumes/wopr-films"] = true
	for _, folder := range c.settings.MediaFolders {
		c.maps[folder] = false
		if _, err := os.Stat(folder); err == nil {
			c.maps[folder] = true
		}
	}

	go c.actor.React()
}

// Stop -
func (c *Core) Stop() {
	mlog.Info("Stopped service Core ...")
}

func (c *Core) getConfig(msg *pubsub.Message) {
	msg.Reply <- &c.settings.Config
	mlog.Info("Sent config")
}

func (c *Core) importMovies(_ *pubsub.Message) {
	t0 := time.Now()
	lib.Notify(c.bus, "import:begin", "Import process started")

	c.wg.Add(1)

	c.bus.Pub(nil, "/command/movie/scan")

	go func() {
		c.wg.Wait()
		lib.Notify(c.bus, "import:end", fmt.Sprintf("Import process finished (%s elapsed)", time.Since(t0).String()))
	}()
}

func (c *Core) addMovie(msg *pubsub.Message) {
	// I'm expecting movie to have the following fields filled
	// Title
	movie := msg.Payload.(*model.Movie)

	// Since it's coming in through this endpoint, we know it's a stub
	movie.FileTitle = movie.Title
	movie.Resolution = cNotAvailable
	movie.FileType = cNotAvailable
	movie.Location = cNotAvailable
	movie.Stub = 1

	// 3 operations, rescrape, update and cache
	c.wg.Add(1)

	// rescrape
	scrape := &pubsub.Message{Payload: movie, Reply: make(chan interface{}, 3)}
	c.bus.Pub(scrape, "/command/movie/scrape")

	go func() {
		c.wg.Wait()
		msg.Reply <- movie
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
}

func (c *Core) doMovieTmdbNotFound(msg *pubsub.Message) {
	item := msg.Payload.(*dto.Scrape)

	store := &pubsub.Message{Payload: item.Movie, Reply: make(chan interface{}, 3)}
	c.bus.Pub(store, "/command/movie/partialstore")
}

func (c *Core) doMovieScraped(msg *pubsub.Message) {
	item := msg.Payload.(*dto.Scrape)

	mlog.Info("ScrapeDTO: %+v", item)

	// I treat the following two commands as one, for the sake of the wg
	// now there are two outstanding locks, which will be decreased by each
	// responding service
	c.wg.Add(1)

	store := &pubsub.Message{Payload: item.Movie, Reply: make(chan interface{}, 3)}
	c.bus.Pub(store, "/command/movie/store")

	cache := &pubsub.Message{Payload: item, Reply: make(chan interface{}, 3)}
	c.bus.Pub(cache, "/command/movie/cache")
}

func (c *Core) doMovieReScraped(msg *pubsub.Message) {
	item := msg.Payload.(*dto.Scrape)

	mlog.Info("ScrapeDTO: %+v", item)

	// I treat the following two commands as one, for the sake of the wg
	// now there are two outstanding locks, which will be decreased by each
	// responding service
	c.wg.Add(1)

	store := &pubsub.Message{Payload: item.Movie, Reply: make(chan interface{}, 3)}
	c.bus.Pub(store, "/command/movie/update")

	cache := &pubsub.Message{Payload: item, Reply: make(chan interface{}, 3)}
	c.bus.Pub(cache, "/command/movie/cache")
}

func (c *Core) pruneMovies(_ *pubsub.Message) {
	t0 := time.Now()

	lib.Notify(c.bus, "prune:begin", "Started Prune Process")

	options := &lib.Options{Offset: 0, Limit: 99999999999999, SortBy: "title", SortOrder: "asc"}
	all := &pubsub.Message{Payload: options, Reply: make(chan interface{}, capacity)}
	c.bus.Pub(all, "/get/movies")

	reply := <-all.Reply
	list := reply.(*model.MoviesDTO)

	if c.settings.UnraidMode {
		hostItems := make(map[string][]*pb.Item)

		for _, host := range c.settings.UnraidHosts {
			hostItems[host] = make([]*pb.Item, 0)
		}

		for _, item := range list.Items {
			index := strings.Index(item.Location, ":")
			if index == -1 {
				// a valid location is wopr:/mnt/user/films/bluray/22 Bullets (2010)/22.Bullets_BLURAY.iso
				// if a ':' isn't found, then this must be a stub
				continue
			}

			host := item.Location[:index]
			location := item.Location[index+1:]

			hostItems[host] = append(hostItems[host], &pb.Item{Id: item.ID, Location: location, Title: item.Title})
		}

		opts := []grpc.DialOption{grpc.WithInsecure()}

		for _, host := range c.settings.UnraidHosts {
			address := fmt.Sprintf("%s.apertoire.org:7624", host)

			conn, err := grpc.Dial(address, opts...)
			if err != nil {
				mlog.Warning("Unable to connect to host (%s): %s", address, err)
				continue
			}
			defer conn.Close()

			client := pb.NewMediaAgentClient(conn)

			rsp, err := client.Exists(context.Background(), &pb.ExistsReq{Items: hostItems[host]})
			if err != nil {
				mlog.Warning("Unable to check exist (%s): %s", address, err)
				continue
			}

			for _, item := range rsp.Items {
				lib.Notify(c.bus, "prune:selected", fmt.Sprintf("UP FOR DELETION: [%d] %s (%s))", item.Id, item.Title, item.Location))

				movie := &pubsub.Message{Payload: item, Reply: make(chan interface{}, capacity)}
				c.bus.Pub(movie, "/command/movie/delete")
			}
		}
	} else {
		for _, item := range list.Items {
			matches := c.re.FindStringSubmatch(item.Location)
			if len(matches) == 0 {
				continue
			}

			folder := filepath.Join("/Volumes", matches[1])
			if !c.maps[folder] {
				continue
			}

			if _, err := os.Stat(item.Location); err != nil {
				if os.IsNotExist(err) {
					lib.Notify(c.bus, "prune:selected", fmt.Sprintf("UP FOR DELETION: [%d] %s (%s))", item.ID, item.Title, item.Location))

					movie := &pubsub.Message{Payload: item, Reply: make(chan interface{}, capacity)}
					c.bus.Pub(movie, "/command/movie/delete")
				}
			}
		}
	}

	lib.Notify(c.bus, "prune:end", fmt.Sprintf("Prune process finished (%s elapsed)", time.Since(t0).String()))
}

func (c *Core) addMediaFolder(msg *pubsub.Message) {
	folder := msg.Payload.(string)
	mlog.Info("addMediaFolder: %s", folder)

	c.settings.MediaFolders = append(c.settings.MediaFolders, folder)
	if err := c.settings.Save(); err != nil {
		mlog.Warning("unable to save settings: %s", err)
	}

	cfg := &pubsub.Message{Payload: c.settings}
	c.bus.Pub(cfg, "/event/config/changed")

	msg.Reply <- &c.settings.Config
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

func (c *Core) doWorkUnitDone(_ *pubsub.Message) {
	c.wg.Done()
}
