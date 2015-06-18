package services

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"jbrodriguez/mediagui/server/lib"
	"jbrodriguez/mediagui/server/model"
	"jbrodriguez/mediagui/server/ws"
)

type Socket struct {
	Service

	bus      *pubsub.PubSub
	settings *lib.Settings
	hub      *ws.Hub

	mailbox chan *pubsub.Mailbox
}

func NewSocket(bus *pubsub.PubSub, settings *lib.Settings) *Socket {
	socket := &Socket{
		bus:      bus,
		settings: settings,
		hub:      ws.NewHub(),
	}

	socket.init()
	return socket
}

func (s *Socket) Start() {
	mlog.Info("Starting service Socket ...")

	s.mailbox = s.register(s.bus, "socket:connections:new", s.doNewConnection)
	s.registerAdditional(s.bus, "socket:connections:broadcast", s.transmit, s.mailbox)

	go s.hub.Run()

	go s.receive()

	go s.react()
}

func (s *Socket) Stop() {
	mlog.Info("Stopped service Socket ...")
}

func (s *Socket) react() {
	for mbox := range s.mailbox {
		// mlog.Info("Socket:Topic: %s", mbox.Topic)
		s.dispatch(mbox.Topic, mbox.Content)
	}
}

func (s *Socket) doNewConnection(msg *pubsub.Message) {
	go s.connect(msg.Payload.(*websocket.Conn))
}

func (s *Socket) connect(wskt *websocket.Conn) {
	c := &ws.Connection{
		Id:   "alpha",
		Ws:   wskt,
		Send: make(chan []byte, 256),
		Hub:  s.hub,
	}

	s.hub.Register <- c

	go c.Writer()
	c.Reader()
}

func (s *Socket) transmit(msg *pubsub.Message) {
	out, err := json.Marshal(msg.Payload)
	if err != nil {
		s.hub.Broadcast <- out
	} else {
		mlog.Warning("Error transmitting websocket message: %s", err)
	}
}

// Receive messages from websocket and dispatch them accordingly
func (s *Socket) receive() {
	for msg := range s.hub.FromWs {
		var wsMsg model.WsMessage

		err := json.Unmarshal(msg, &wsMsg)
		if err != nil {
			mlog.Warning("Unable to unmarshal: %s", err)
		}

		pkt := &pubsub.Message{Payload: wsMsg}
		s.bus.Pub(pkt, wsMsg.Topic)
	}
}
