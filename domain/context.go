package domain

import "github.com/cskr/pubsub"

type Context struct {
	MediaFolders []string       `json:"mediaFolders"`
	UnraidHosts  []string       `json:"unraidHosts"`
	Version      string         `json:"version"`
	DataDir      string         `json:"-"`
	Hub          *pubsub.PubSub `json:"-"`
	UserAgent    string         `json:"-"`
	TmdbKey      string         `json:"-"`
}
