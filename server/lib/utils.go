package lib

import (
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/nfnt/resize"

	"mediagui/dto"
)

// Exists - Check if File / Directory Exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// SearchFile - look for file in the specified locations
func SearchFile(name string, locations []string) string {
	for _, location := range locations {
		if b, _ := Exists(filepath.Join(location, name)); b {
			return location
		}
	}

	return ""
}

// RestGet -
func RestGet(url string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/604.1.22 (KHTML, like Gecko) Version/10.2 Safari/604.1.22")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer Close(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Download -
func Download(url, dst string) (err error) {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return err
}

// Notify -
func Notify(bus *pubsub.PubSub, topic, text string) {
	mlog.Info(text)
	payload := &dto.Packet{Topic: topic, Payload: text}
	bus.Pub(&pubsub.Message{Payload: payload}, "socket:broadcast")
}

// ResizeImage -
func ResizeImage(src, dst string) (err error) {
	// open "test.jpg"
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}
	// file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(80, 0, img, resize.Lanczos3)

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// write new image to file
	return jpeg.Encode(out, m, nil)
}

// Close -
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
}
