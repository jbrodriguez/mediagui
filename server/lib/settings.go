package lib

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/namsral/flag"
)

// Config -
type Config struct {
	UnraidMode   bool     `json:"unraidMode"`
	UnraidHosts  []string `json:"unraidHosts"`
	MediaFolders []string `json:"mediaFolders"`
	Version      string   `json:"version"`
}

// Settings -
type Settings struct {
	Config

	DataDir string
	WebDir  string
	LogDir  string

	Location   string
	GinMode    string
	CPUProfile string
	CheckOnly  bool
	WorkDir    string
}

// NewSettings -
func NewSettings(name, version, home string, locations []string) (*Settings, error) {
	var config, dataDir, webDir, logDir, mediaFolders, ginMode, cpuprofile, unraidHosts, workDir string
	var logtostderr, unraidMode, checkOnly bool

	location := SearchFile(name, locations)

	flag.BoolVar(&logtostderr, "logtostderr", true, "true/false log to stderr")
	flag.StringVar(&config, "config", "", "config location")
	flag.StringVar(&dataDir, "datadir", filepath.Join(home, ".mediagui/db"), "folder containing the database files")
	flag.StringVar(&webDir, "webdir", filepath.Join(home, ".mediagui/web"), "folder where web app will be read from")
	flag.StringVar(&logDir, "logdir", "", "folder where log file will be written to")
	flag.StringVar(&mediaFolders, "mediafolders", "/mnt/user/films", "folders that will be scanned for media")
	flag.StringVar(&ginMode, "gin_mode", "release", "gin mode")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to file")
	flag.BoolVar(&unraidMode, "unraid_mode", true, "if true the app will work distributed with a service running on the unraid host")
	flag.StringVar(&unraidHosts, unraidHosts, "wopr|hal", "specify which unraid hosts will be scanned for movies. the service agent must be running in that host")
	flag.BoolVar(&checkOnly, "checkOnly", false, "for importer cli tool -> true: check if exists (by name) / false: add movies as stub")
	flag.StringVar(&workDir, "workDir", filepath.Join(home, "tmp", "mediagui"), "for importer cli tool -> folder where the file with movies to import resides")

	if found, _ := Exists(location); found {
		if err := flag.Set("config", location); err != nil {
			fmt.Println("unable to set config file")
		}
	}
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
	s.Location = location
	s.GinMode = ginMode
	s.CPUProfile = cpuprofile
	s.UnraidMode = unraidMode
	s.CheckOnly = checkOnly
	s.WorkDir = workDir
	if unraidHosts == "" {
		s.UnraidHosts = make([]string, 0)
	} else {
		s.UnraidHosts = strings.Split(unraidHosts, "|")
	}

	return s, nil
}

// Save -
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

	return writeLine(file, "mediafolders", mediaFolders)
}

func writeLine(file *os.File, key, value string) error {
	_, err := io.WriteString(file, fmt.Sprintf("%s=%s\n", key, value))
	if err != nil {
		return err
	}

	return nil
}
