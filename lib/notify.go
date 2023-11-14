package lib

import (
	"github.com/cskr/pubsub"

	"mediagui/domain"
	"mediagui/logger"
)

func Notify(hub *pubsub.PubSub, topic, text string) {
	logger.Blue(text)
	payload := &domain.Packet{Topic: topic, Payload: text}
	hub.Pub(payload, "socket:broadcast")
}
