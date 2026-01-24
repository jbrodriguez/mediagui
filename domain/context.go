package domain

import (
	"time"

	"github.com/cskr/pubsub"
)

type Context struct {
	MediaFolders []string       `json:"mediaFolders"`
	UnraidHosts  []string       `json:"unraidHosts"`
	Version      string         `json:"version"`
	DataDir      string         `json:"-"`
	Hub          *pubsub.PubSub `json:"-"`
	UserAgent    string         `json:"-"`
	TmdbKey      string         `json:"-"`
	OmdbKey      string         `json:"-"`

	// Rate limiting configuration
	TmdbRateLimit    int           `json:"-"`
	TmdbRateInterval time.Duration `json:"-"`
	ImdbRateLimit    int           `json:"-"`
	ImdbRateInterval time.Duration `json:"-"`
}
