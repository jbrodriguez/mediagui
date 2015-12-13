package main

import (
	log "github.com/golang/glog"
	"github.com/myodc/go-micro/cmd"
	// "github.com/myodc/go-micro/examples/server/handler"
	// "github.com/myodc/go-micro/examples/server/subscriber"
	"github.com/myodc/go-micro/server"
	"os"
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
		server.Name("io.jbrodriguez.mediagui.scanner."+host),
		server.Address("0.0.0.0:0"),
	)

	// Register Handlers
	server.Handle(
		server.NewHandler(
			&Scanner{host: host},
		),
	)

	log.Info("Scanner started ")

	// Run server
	if err := server.Run(); err != nil {
		log.Info(err)
	}

}
