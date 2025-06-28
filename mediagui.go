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
	MediaFolders []string
	UnraidHosts  []string
	UserAgent    string
	TmdbKey      string

	Boot cmd.Boot `cmd:"" default:"1" help:"start processing"`
}

func main() {
	home := os.Getenv("HOME")
	dataDir := filepath.Join(home, ".local", "share", "mediagui")

	ctx := kong.Parse(&cli)
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
