package lib

import (
	"encoding/json"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"jbrodriguez/mediagui/server/model"
	"net/http"
	"os"
)

// Check if File / Directory Exists
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

func Notify(bus *pubsub.PubSub, topic, text string) {
	mlog.Info(text)
	payload := &model.Packet{Topic: topic, Payload: text}
	bus.Pub(&pubsub.Message{Payload: payload}, "socket:connections:broadcast")
}
