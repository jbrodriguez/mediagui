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

	if settings.LogDir != "" {
		mlog.Start(mlog.LevelInfo, filepath.Join(settings.LogDir, "mediagui.log"))
	} else {
		mlog.Start(mlog.LevelInfo, "")
	}

	mlog.Info("mediagui v%s starting ...", Version)

	bus := pubsub.New(623)

	// socket := services.NewSocket(bus, &settings)
	server := services.NewServer(bus, settings)
	core := services.NewCore(bus, settings)
	dal := services.NewDal(bus, settings)

	// socket.Start()
	dal.Start()
	server.Start()
	core.Start()

	mlog.Info("Press Ctrl+C to stop ...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		mlog.Info("Received an interrupt, shutting the app down ...")

		core.Stop()
		server.Stop()
		dal.Stop()
		// socket.Stop()

		break
	}

	mlog.Stop()
}
