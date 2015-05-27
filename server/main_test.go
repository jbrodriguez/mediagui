package main

import (
	"fmt"
	"github.com/jbrodriguez/mlog"
	// "github.com/jbrodriguez/pubsub"
	"github.com/stretchr/testify/assert"
	"io"
	"jbrodriguez/mediagui/server/lib"
	// "jbrodriguez/mediagui/server/services"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(filename, text string) error {
	file, err := os.Create(filename)
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = io.WriteString(file, fmt.Sprintf(text))
	if err != nil {
		return err
	}

	return nil
}

func TestSettingsNotFound(t *testing.T) {
	home := os.Getenv("HOME")
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	log.Printf("Unable to get current directory: %s", err.Error)
	// 	os.Exit(1)
	// }

	locations := []string{
		filepath.Join(home, "tmp/mgtest/.mediagui/mediagui.conf"),
		filepath.Join(home, "tmp/mgtest/usr/local/etc/mediagui.conf"),
		filepath.Join(home, "tmp/mgtest/mediagui.conf"),
	}

	settings, err := lib.NewSettings(Version, home, locations)

	if assert.Error(t, err) {
		mlog.Info("Ok: error was:\n%s", err.Error())
	} else {
		mlog.Fatalf("settings: %+v", settings)
	}
}

func TestSettingsFound(t *testing.T) {
	home := os.Getenv("HOME")

	// cwd, err := os.Getwd()
	// if err != nil {
	// 	mlog.Info("Unable to get current directory: %s", err.Error)
	// 	os.Exit(1)
	// }
	path := filepath.Join(home, "tmp/mgtest/.mediagui")
	os.MkdirAll(path, 0777)
	defer os.RemoveAll(path)

	text := "datadir=mg_datadir\nwebdir=mg_webdir\nmediafolders=movies/bluray|tv shows|movies/blurip"
	err := write(filepath.Join(path, "mediagui.conf"), text)

	assert.NoError(t, err)

	b, err := lib.Exists(filepath.Join(path, "mediagui.conf"))

	assert.Equal(t, true, b)

	locations := []string{
		filepath.Join(home, "tmp/mgtest/.mediagui/mediagui.conf"),
		filepath.Join(home, "tmp/mgtest/usr/local/etc/mediagui.conf"),
		filepath.Join(home, "tmp/mgtest/mediagui.conf"),
	}

	settings, err := lib.NewSettings(Version, home, locations)

	if assert.NoError(t, err) {
		assert.Equal(t, "mg_datadir", settings.DataDir)
		assert.Equal(t, "mg_webdir", settings.WebDir)
		assert.Equal(t, "", settings.LogDir)
		assert.Equal(t, strings.Split("movies/bluray|tv shows|movies/blurip", "|"), settings.MediaFolders)
	}

}

// func TestServerService(t *testing.T) {
// 	home := os.Getenv("HOME")

// 	// cwd, err := os.Getwd()
// 	// if err != nil {
// 	// 	mlog.Info("Unable to get current directory: %s", err.Error)
// 	// 	os.Exit(1)
// 	// }
// 	path := filepath.Join(home, "tmp/mgtest/.mediagui")
// 	os.MkdirAll(path, 0777)
// 	defer os.RemoveAll(path)

// 	html := filepath.Join(home, "tmp/mgtest/web/index.html")
// 	text := fmt.Sprintf("datadir=mg_datadir\nwebdir=%s\nmediafolders=movies/bluray|tv shows|movies/blurip", html)
// 	err := write(filepath.Join(path, "mediagui.conf"), text)

// 	assert.NoError(t, err)

// 	b, err := lib.Exists(filepath.Join(path, "mediagui.conf"))

// 	assert.Equal(t, true, b)

// 	locations := []string{
// 		filepath.Join(home, "tmp/mgtest/.mediagui/mediagui.conf"),
// 		filepath.Join(home, "tmp/mgtest/usr/local/etc/mediagui.conf"),
// 		filepath.Join(home, "tmp/mgtest/mediagui.conf"),
// 	}

// 	settings, err := lib.NewSettings(Version, home, locations)

// 	if assert.NoError(t, err) {
// 		assert.Equal(t, "mg_datadir", settings.DataDir)
// 		assert.Equal(t, "mg_webdir", settings.WebDir)
// 		assert.Equal(t, "", settings.LogDir)
// 		assert.Equal(t, strings.Split("movies/bluray|tv shows|movies/blurip", "|"), settings.MediaFolders)
// 	}

// 	text = "<body><h1>we shall overcome</h1></body>"
// 	err = write(html, text)
// 	assert.NoError(t, err)

// 	bus := pubsub.New(623)

// 	server := services.NewServer(bus, settings)
// 	server.Start()
// }

func TestMain(m *testing.M) {
	mlog.Start(mlog.LevelInfo, "")

	home := os.Getenv("HOME")
	path := filepath.Join(home, "tmp/mgtest")
	os.RemoveAll(path)

	ret := m.Run()

	os.RemoveAll(path)

	// mlog.Stop()

	os.Exit(ret)
}
