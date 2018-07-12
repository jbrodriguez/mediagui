package agent

import (
	"jbrodriguez/mediagui/server/src/proto"

	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"
)

// Agent -
type Agent struct {
	Host string
}

// Scan -
func (s *Agent) Scan(ctx context.Context, req *agent.ScanReq, rsp *agent.ScanRsp) error {
	log.Printf("Received Agent.Scan request: %v\n", req)

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

	log.Printf("Sent back %d files\n", len(rsp.Filenames))

	return nil
}

func (s *Agent) walk(folder, mask string) []string {
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

		files = append(files, s.Host+":"+path)

		return nil
	})

	return files
}

// Exists -
func (s *Agent) Exists(ctx context.Context, req *agent.ExistsReq, rsp *agent.ExistsRsp) error {
	log.Printf("Received Agent.Exists request: %d items", len(req.Items))

	rsp.Items = make([]*agent.Item, 0)

	for _, item := range req.Items {
		exists := true

		if _, err := os.Stat(item.Location); err != nil {
			exists = !os.IsNotExist(err)
		}

		if !exists {
			log.Printf("Location %s doesn't exist\n", item.Location)
			rsp.Items = append(rsp.Items, item)
		}
	}

	return nil
}
