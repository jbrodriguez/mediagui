package main

import (
	"mediagui/domain"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
	"github.com/cskr/pubsub"

	"mediagui/cmd"
)

var Version string

var cli struct {
	MediaFolders []string `env:"MEDIA_FOLDERS" help:"Media folders to scan"`
	UnraidHosts  []string `env:"UNRAID_HOSTS" help:"Unraid hosts to monitor"`
	UserAgent    string   `env:"USER_AGENT" help:"User agent for API requests"`
	TmdbKey      string   `env:"TMDB_KEY" help:"TMDB API key"`
	OmdbKey      string   `env:"OMDB_KEY" help:"OMDb API key"`
	DataDir      string   `env:"DATA_DIR" help:"Data directory for application files"`

	Boot cmd.Boot `cmd:"" default:"1" help:"start processing"`
}

func main() {
	ctx := kong.Parse(&cli)

	// Set default data directory if not provided
	dataDir := cli.DataDir
	if dataDir == "" {
		home := os.Getenv("HOME")
		dataDir = filepath.Join(home, ".local", "share", "mediagui")
	}

	err := ctx.Run(&domain.Context{
		DataDir:      dataDir,
		MediaFolders: cli.MediaFolders,
		UnraidHosts:  cli.UnraidHosts,
		UserAgent:    cli.UserAgent,
		TmdbKey:      cli.TmdbKey,
		OmdbKey:      cli.OmdbKey,
		Version:      Version,
		Hub:          pubsub.New(23),

		// TMDB rate limiting: 40 requests per 10 seconds (official limit)
		TmdbRateLimit:    40,
		TmdbRateInterval: 250 * time.Millisecond, // 40 req/10s = 1 req per 250ms

		// IMDB rate limiting: more conservative for web scraping
		ImdbRateLimit:    10,
		ImdbRateInterval: 1 * time.Second, // 10 req per 10s = 1 req per second
	})
	ctx.FatalIfErrorf(err)
}
