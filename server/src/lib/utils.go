package lib

import (
	"encoding/json"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
	"jbrodriguez/mediagui/server/src/dto"
	"net/http"
	"os"
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

// RestGet -
func RestGet(url string, reply interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	err = json.NewDecoder(resp.Body).Decode(reply)

	return err
}

// Download -
func Download(url, dst string) (err error) {
	out, err := os.Create(dst)
	if err != nil {
		// mlog.Info("Unable to create: %s", dst)
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		// mlog.Info("Unable to download %s", url)
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		// mlog.Info("unable to save to %s", dst)
		return err
	}

	return err
}

// Notify -
func Notify(bus *pubsub.PubSub, topic, text string) {
	mlog.Info(text)
	payload := &dto.Packet{Topic: topic, Payload: text}
	// payload := &dto.Packet{Topic: topic, Payload: text}
	bus.Pub(&pubsub.Message{Payload: payload}, "socket:connections:broadcast")
}

// ResizeImage -
func ResizeImage(src, dst string) (err error) {
	// open "test.jpg"
	file, err := os.Open(src)
	if err != nil {
		// mlog.Error(err)
		return err
	}
	defer file.Close()

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		// mlog.Error(err)
		return err
	}
	// file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(80, 0, img, resize.Lanczos3)

	out, err := os.Create(dst)
	if err != nil {
		// mlog.Error(err)
		return err
	}
	defer out.Close()

	// write new image to file
	return jpeg.Encode(out, m, nil)
}
