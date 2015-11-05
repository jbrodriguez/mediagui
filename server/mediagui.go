package main

import (
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"jbrodriguez/mediagui/server/lib"
	"jbrodriguez/mediagui/server/services"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
)

var Version string

func main() {
	// look for mediagui.conf at the following places
	// $HOME/.mediagui/mediagui.conf
	// /usr/local/etc/mediagui.conf
	// <current dir>/mediagui.conf
	home := os.Getenv("HOME")

	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Unable to get current directory: %s", err.Error())
		os.Exit(1)
	}

	locations := []string{
		filepath.Join(home, ".mediagui/mediagui.conf"),
		"/usr/local/etc/mediagui.conf",
		filepath.Join(cwd, "mediagui.conf"),
	}

	settings, err := lib.NewSettings(Version, home, locations)
	if err != nil {
		log.Printf("Unable to start the app: %s", err.Error())
		os.Exit(2)
	}

	if settings.CpuProfile != "" {
		f, err := os.Create(settings.CpuProfile)
		if err != nil {
			log.Printf("Unable to set up profiling: %s", err)
			os.Exit(3)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if settings.LogDir != "" {
		mlog.Start(mlog.LevelInfo, filepath.Join(settings.LogDir, "mediagui.log"))
	} else {
		mlog.Start(mlog.LevelInfo, "")
	}

	mlog.Info("mediagui v%s starting ...", Version)
	mlog.Info("using config file located at %s", settings.Location)

	bus := pubsub.New(8623)

	dal := services.NewDal(bus, settings)
	socket := services.NewSocket(bus, settings)
	server := services.NewServer(bus, settings)
	scanner := services.NewScanner(bus, settings)
	scraper := services.NewScraper(bus, settings)
	cache := services.NewCache(bus, settings)
	core := services.NewCore(bus, settings)

	dal.Start()
	socket.Start()
	server.Start()
	scanner.Start()
	scraper.Start()
	cache.Start()
	core.Start()

	mlog.Info("Press Ctrl+C to stop ...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		mlog.Info("Received an interrupt, shutting the app down ...")

		core.Stop()
		cache.Stop()
		scraper.Stop()
		scanner.Stop()
		server.Stop()
		socket.Stop()
		dal.Stop()

		break
	}

	mlog.Stop()
}
