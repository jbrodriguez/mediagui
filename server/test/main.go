package main

import (
	"jbrodriguez/mediagui/server/src/proto"

	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"

	"golang.org/x/net/context"
)

func main() {
	cmd.Init()

	folders := []string{
		`/mnt/user/films`,
	}

	// Create new request to service go.micro.srv.example, method Example.Call
	req := client.NewRequest("io.jbrodriguez.mediagui.agent", "Scanner.Scan", &agent.ScanReq{
		// Folders: s.settings.MediaFolders,
		Folders: folders,
	})

	fmt.Printf("req=%+v\n", req)

	rsp := &agent.ScanRsp{}

	t0 := time.Now()

	// Call service
	if err := client.Call(context.Background(), req, rsp); err != nil {
		fmt.Printf("WARNING: Unable to connect to scanning service: %s\n", err)
		return
	}

	𝛥t := float64(time.Since(t0)) / 1e9

	for _, file := range rsp.Filenames {
		fmt.Printf("file=%s\n", file)
	}

	fmt.Printf("walked %d files in %.3f seconds\n", len(rsp.Filenames), 𝛥t)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	fmt.Printf("Received signal %s", <-ch)
}
