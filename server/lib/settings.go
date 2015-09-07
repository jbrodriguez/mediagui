package lib

import (
	"errors"
	"fmt"
	"github.com/namsral/flag"
	"io"
	"os"
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

	Location   string
	GinMode    string
	CpuProfile string
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

	var config, dataDir, webDir, logDir, mediaFolders, ginMode, cpuprofile string
	flag.StringVar(&config, "config", "", "config location")
	flag.StringVar(&dataDir, "datadir", filepath.Join(home, ".mediagui/db"), "folder containing the database files")
	flag.StringVar(&webDir, "webdir", filepath.Join(home, ".mediagui/web"), "folder where web app will be read from")
	flag.StringVar(&logDir, "logdir", "", "folder where log file will be written to")
	flag.StringVar(&mediaFolders, "mediafolders", "", "folders that will be scanned for media")
	flag.StringVar(&ginMode, "gin_mode", "release", "gin mode")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to file")

	flag.Set("config", location)
	flag.Parse()

	// fmt.Printf("mediaFolders: %s\n", mediaFolders)

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
	s.Location = location
	s.GinMode = ginMode
	s.CpuProfile = cpuprofile

	return s, nil
}

func (s *Settings) Save() error {
	file, err := os.Create(s.Location)
	defer file.Close()

	if err != nil {
		return err
	}

	if err = writeLine(file, "datadir", s.DataDir); err != nil {
		return err
	}

	if err = writeLine(file, "webdir", s.WebDir); err != nil {
		return err
	}

	if err = writeLine(file, "logdir", s.LogDir); err != nil {
		return err
	}

	mediaFolders := strings.Join(s.MediaFolders, "|")
	if err = writeLine(file, "mediafolders", mediaFolders); err != nil {
		return err
	}

	return nil
}

func writeLine(file *os.File, key, value string) error {
	_, err := io.WriteString(file, fmt.Sprintf("%s=%s\n", key, value))
	if err != nil {
		return err
	}

	return nil
}
