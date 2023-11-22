package domain

import "github.com/cskr/pubsub"

type Context struct {
	MediaFolder string         `json:"mediaFolder"`
	UnraidHosts []string       `json:"unraidHosts"`
	Version     string         `json:"version"`
	DataDir     string         `json:"-"`
	Hub         *pubsub.PubSub `json:"-"`
	UserAgent   string         `json:"-"`
	TmdbKey     string         `json:"-"`
}
