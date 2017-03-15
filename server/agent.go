// +build agent

package main

import (
	"jbrodriguez/mediagui/server/src/agent"

	"log"
	"os"

	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"
)

func main() {
	// optionally setup command line usage
	cmd.Init()

	host, err := os.Hostname()
	if err != nil {
		log.Fatalf("Unable to obtain hostname: %s", err)
	}

	// Initialise Server
	server.Init(
		server.Name("io.jbrodriguez.mediagui.agent."+host),
		server.Address("0.0.0.0:0"),
	)

	// Register Handlers
	server.Handle(
		server.NewHandler(
			&agent.Agent{Host: host},
		),
	)

	// Run server
	if err := server.Run(); err != nil {
		log.Printf("%s", err)
	}

}
