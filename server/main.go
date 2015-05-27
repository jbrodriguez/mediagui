package main

import (
	"github.com/jbrodriguez/mlog"
	"github.com/julienschmidt/httprouter"
	"jbrodriguez/mediagui/server/lib"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var Version string

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/Volumes/Users/kayak/code/src/jbrodriguez/mediagui/target/build/index.html")
}

func main() {
	// look for mediagui.conf at the following places
	// $HOME/.mediagui/mediagui.conf
	// /usr/local/etc/mediagui.conf
	// <current dir>/mediagui.conf
	home := os.Getenv("HOME")

	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Unable to get current directory: %s", err.Error())
		os.Exit(1)
	}

	locations := []string{
		filepath.Join(home, ".mediagui/mediagui.conf"),
		"/usr/local/etc/mediagui.conf",
		filepath.Join(cwd, "mediagui.conf"),
	}

	settings, err := lib.NewSettings(Version, home, locations)
	if err != nil {
		log.Printf("Unable to start the app: %s", err.Error())
		os.Exit(2)
	}

	mlog.Start(mlog.LevelInfo, settings.LogDir)

	mlog.Info("mediagui v%s starting ...", Version)

	router := httprouter.New()

	router.HandlerFunc("GET", "/", index)
	router.ServeFiles("/app/*filepath", http.Dir("/Volumes/Users/kayak/code/src/jbrodriguez/mediagui/target/build//app"))

	log.Fatal(http.ListenAndServe(":7623", router))
}
