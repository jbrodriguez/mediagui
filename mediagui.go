package main

import (
	"mediagui/domain"
	"os"
	"path/filepath"

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
		Version:      Version,
		Hub:          pubsub.New(23),
	})
	ctx.FatalIfErrorf(err)
}
