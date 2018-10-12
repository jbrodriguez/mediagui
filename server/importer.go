// +build importer

package main

import (
	"log"
	"os"

	"mediagui/importer"
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
