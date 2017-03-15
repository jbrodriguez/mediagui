package main

import (
	"jbrodriguez/mediagui/proto"

	"os"
	"path/filepath"
	"strings"

	log "github.com/golang/glog"

	"golang.org/x/net/context"
)

// Agent -
type Agent struct {
	host string
}

// Scan -
func (s *Agent) Scan(ctx context.Context, req *agent.ScanReq, rsp *agent.ScanRsp) error {
	log.Infof("Received Agent.Scan request: %v", req)

	// var files []string

	for _, folder := range req.Folders {
		list := s.walk(folder, req.Mask)
		// if err != nil {
		// 	mlog.Warning("Unable to scan folder (%s): %s", folder, err)
		// }

		// files = append(files, list...)
		rsp.Filenames = append(rsp.Filenames, list...)
	}

	// for _, f := range files {
	// 	rsp.Filenames = append(rsp.Filenames, f)
	// }

	log.Infof("Sent back %d files", len(rsp.Filenames))

	return nil
}

func (s *Agent) walk(folder, mask string) []string {
	if folder[len(folder)-1] != '/' {
		folder = folder + "/"
	}

	var files []string

	filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			log.Infof("Agent.Scan.walk: %s (%s) - [%+v]", err, path, f)
		}

		if f.IsDir() {
			return nil
		}

		if !strings.Contains(mask, strings.ToLower(filepath.Ext(path))) {
			// mlog.Info("[%s] excluding %s", filepath.Ext(path), path)
			return nil
		}

		// log.Infof("file=%s", path)

		files = append(files, s.host+":"+path)

		return nil
	})

	return files
}

// Exists -
func (s *Agent) Exists(ctx context.Context, req *agent.ExistsReq, rsp *agent.ExistsRsp) error {
	// log.Infof("Received Agent.Exists request: %v", req)

	rsp.Exists = true
	if _, err := os.Stat(req.Location); err != nil {
		rsp.Exists = !os.IsNotExist(err)
	}

	if !rsp.Exists {
		log.Infof("Location %s doesn't exist", req.Location)
	}

	return nil
}
