package app

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"

	"mediagui/lib"
	"mediagui/services"
)

// App empty placeholder
type App struct {
}

// Setup app
func (a *App) Setup(version string) (*lib.Settings, error) {
	// look for mediagui.conf at the following places
	// /boot/config/plugins/mediagui/
	// <current dir>/mediagui.conf
	home := os.Getenv("HOME")

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	locations := []string{
		filepath.Join(home, ".mediagui/mediagui.conf"),
		"/usr/local/etc/mediagui.conf",
		filepath.Join(cwd, "mediagui.conf"),
	}

	settings, err := lib.NewSettings("mediagui.conf", version, home, locations)

	return settings, err
}

// Run app
func (a *App) Run(settings *lib.Settings) {
	if settings.LogDir != "" {
		mlog.Start(mlog.LevelInfo, filepath.Join(settings.LogDir, "mediagui.log"))
	} else {
		mlog.Start(mlog.LevelInfo, "")
	}

	mlog.Info("mediagui v%s starting ...", settings.Version)

	var msg string
	if settings.Location == "" {
		msg = "No config file specified. Using app defaults ..."
	} else {
		msg = fmt.Sprintf("Using config file located at %s ...", settings.Location)
	}
	mlog.Info(msg)

	bus := pubsub.New(8623)

	dal := services.NewDal(bus, settings)
	server := services.NewServer(bus, settings)
	scanner := services.NewScanner(bus, settings)
	scraper := services.NewScraper(bus, settings)
	cache := services.NewCache(bus, settings)
	core := services.NewCore(bus, settings)

	dal.Start()
	server.Start()
	scanner.Start()
	scraper.Start()
	cache.Start()
	core.Start()

	mlog.Info("Press Ctrl+C to stop ...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	mlog.Info("Received signal: (%s) ... shutting down the app now ...", <-c)

	core.Stop()
	cache.Stop()
	scraper.Stop()
	scanner.Stop()
	server.Stop()
	dal.Stop()

	if err := mlog.Stop(); err != nil {
		os.Exit(2)
	}
}
