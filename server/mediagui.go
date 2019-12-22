package main

import (
	"log"
	"os"

	"mediagui/app"
)

// Version - plugin version
var Version string

func main() {
	mg := app.App{}

	settings, err := mg.Setup(Version)
	if err != nil {
		log.Printf("Unable to start the app: %s", err)
		os.Exit(1)
	}

	mg.Run(settings)
}
