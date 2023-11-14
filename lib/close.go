package lib

import (
	"io"
	"log"
)

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
}
