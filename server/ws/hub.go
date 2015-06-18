package ws

import ()

const (
	// pongWait   = 60 * time.Second
	bufferSize = 8192
)

type Hub struct {
	// registered connections
	connections map[*Connection]bool

	// inbound messages from connections
	Broadcast chan []byte

	// inbound messages from connections
	Emit chan []byte

	FromWs chan []byte

	// register requests from connection
	Register chan *Connection

	// unregister request from connection
	Unregister chan *Connection
}

func NewHub() *Hub {
	hub := &Hub{}

	hub.connections = make(map[*Connection]bool)
	hub.Register = make(chan *Connection)
	hub.Unregister = make(chan *Connection)
	hub.Broadcast = make(chan []byte)
	hub.Emit = make(chan []byte)
	hub.FromWs = make(chan []byte, 256)
	// hub.broadcast = bus.Sub("socket:broadcast")
	// hub.emit = bus.Sub("socket:emit")

	return hub
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.connections[c] = true
		case c := <-h.Unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.Send)
			}
		case m := <-h.Broadcast:
			for c := range h.connections {
				select {
				case c.Send <- m:
				default:
					close(c.Send)
					delete(h.connections, c)
				}
			}
		}
	}
}
