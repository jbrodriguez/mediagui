package lib

import (
	"errors"
	"github.com/namsral/flag"
	"path/filepath"
	"strings"
)

type Config struct {
	MediaFolders []string `json:"mediaFolders"`
	Version      string   `json:"version"`
}

type Settings struct {
	Config

	DataDir string
	WebDir  string
	LogDir  string
}

func searchConfig(locations []string) string {
	for _, location := range locations {
		if b, _ := Exists(location); b {
			return location
		}
	}

	return ""
}

func NewSettings(version, home string, locations []string) (*Settings, error) {
	location := searchConfig(locations)
	if location == "" {
		msg := "Unable to find mediagui.conf\nIt should be placed at any these locations:\n$HOME/.mediagui/\n/usr/local/etc\n<app directory>"
		return nil, errors.New(msg)
	}

	var config, dataDir, webDir, logDir, mediaFolders string
	flag.StringVar(&config, "config", "", "config location")
	flag.StringVar(&dataDir, "datadir", filepath.Join(home, ".mediagui/data"), "folder containing the database files")
	flag.StringVar(&webDir, "webdir", filepath.Join(home, ".mediagui/web"), "folder where web app will be read from")
	flag.StringVar(&logDir, "logdir", "", "folder where log file will be written to")
	flag.StringVar(&mediaFolders, "mediafolders", "", "folders that will be scanned for media")

	flag.Set("config", location)
	flag.Parse()

	s := &Settings{}
	if mediaFolders == "" {
		s.MediaFolders = make([]string, 0)
	} else {
		s.MediaFolders = strings.Split(mediaFolders, "|")
	}
	s.Version = version
	s.DataDir = dataDir
	s.WebDir = webDir
	s.LogDir = logDir
	return s, nil
}
