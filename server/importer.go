// +build importer

package main

import (
	"jbrodriguez/mediagui/server/src/importer"
	"log"
	"os"
)

// Version - plugin version
var Version string

func main() {
	importer := importer.App{}

	settings, err := importer.Setup(Version)
	if err != nil {
		log.Printf("Unable to start the importer: %s", err)
		os.Exit(1)
	}

	importer.Run(settings)
}
