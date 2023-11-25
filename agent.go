//go:build agent
// +build agent

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"

	pb "mediagui/mediaagent"
)

// var Version string
// Version - plugin version
var (
	Version string
	port    = flag.Int("port", 7624, "The server port")
)

type mediaAgentServer struct {
	host string
	pb.UnimplementedMediaAgentServer
}

// Scan -
func (m *mediaAgentServer) Scan(_ context.Context, req *pb.ScanReq) (*pb.ScanRsp, error) {
	log.Printf("Received Agent.Scan request: %v\n", req)

	rep := &pb.ScanRsp{
		Filenames: make([]string, 0),
	}

	for _, folder := range req.Folders {
		list := m.walk(folder, req.Mask)
		// if err != nil {
		// 	mlog.Warning("Unable to scan folder (%s): %s", folder, err)
		// }

		// files = append(files, list...)
		rep.Filenames = append(rep.Filenames, list...)
	}

	// for _, f := range files {
	// 	rsp.Filenames = append(rsp.Filenames, f)
	// }

	log.Printf("Sent back %d files\n", len(rep.Filenames))

	return rep, nil
}

func (m *mediaAgentServer) walk(folder, mask string) []string {
	if folder[len(folder)-1] != '/' {
		folder = folder + "/"
	}

	var files []string

	filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Agent.Scan.walk: %s (%s) - [%+v]\n", err, path, f)
		}

		if f.IsDir() {
			return nil
		}

		if !strings.Contains(mask, strings.ToLower(filepath.Ext(path))) {
			// mlog.Info("[%s] excluding %s", filepath.Ext(path), path)
			return nil
		}

		// log.Infof("file=%s", path)

		files = append(files, m.host+":"+path)

		return nil
	})

	return files
}

// Exists -
func (m *mediaAgentServer) Exists(_ context.Context, req *pb.ExistsReq) (*pb.ExistsRsp, error) {
	log.Printf("Received Agent.Exists request: %d items", len(req.Items))

	rep := &pb.ExistsRsp{
		Items: make([]*pb.Item, 0),
	}

	for _, item := range req.Items {
		exists := true

		if _, err := os.Stat(item.Location); err != nil {
			exists = !os.IsNotExist(err)
		}

		if !exists {
			log.Printf("Location %s doesn't exist\n", item.Location)
			rep.Items = append(rep.Items, item)
		}
	}

	return rep, nil
}

func newServer() *mediaAgentServer {
	host, err := os.Hostname()
	if err != nil {
		log.Fatalf("Unable to obtain hostname: %s", err)
	}

	return &mediaAgentServer{host: host}
}

func main() {
	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%d", *port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMediaAgentServer(grpcServer, newServer())
	log.Printf("started mediaagent v%s listening on %s", Version, address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
