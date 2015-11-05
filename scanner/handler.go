package main

import (
	// "github.com/MichaelTJones/walk"
	log "github.com/golang/glog"
	"jbrodriguez/mediagui/proto"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"
)

const includedMask = ".bdmv|.iso|.img|.nrg|.mkv|.avi|.xvid|.ts|.mpg|.dvr-ms|.mdf|.wmv|.ifo"

type Scanner struct {
	host string
}

func (s *Scanner) Scan(ctx context.Context, req *scan.Request, rsp *scan.Response) error {
	log.Infof("Received scanner.Scan request: %v", req)

	// var files []string

	for _, folder := range req.Folders {
		list := s.walk(folder)
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

func (s *Scanner) walk(folder string) []string {
	if folder[len(folder)-1] != '/' {
		folder = folder + "/"
	}

	var files []string

	filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			log.Infof("scanner.walk: %s (%s) - [%+v]", err, path, f)
		}

		if f.IsDir() {
			return nil
		}

		if !strings.Contains(includedMask, strings.ToLower(filepath.Ext(path))) {
			// mlog.Info("[%s] excluding %s", filepath.Ext(path), path)
			return nil
		}

		// log.Infof("file=%s", path)

		files = append(files, s.host+":"+path)

		return nil
	})

	return files
}
