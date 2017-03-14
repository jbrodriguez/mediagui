package main

import (
	"jbrodriguez/mediagui/server/src/app"
	"log"
	"os"
)

// Version - plugin version
var Version string

func main() {
	app := app.App{}

	settings, err := app.Setup(Version)
	if err != nil {
		log.Printf("Unable to start the app: %s", err)
		os.Exit(1)
	}

	app.Run(settings)
}
