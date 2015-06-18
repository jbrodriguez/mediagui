package ws

import (
	// "apertoire.net/unbalance/server/dto"
	"github.com/gorilla/websocket"
	"github.com/jbrodriguez/mlog"
	// "github.com/jbrodriguez/pubsub"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Connection struct {
	Id   string
	Ws   *websocket.Conn
	Send chan []byte
	Hub  *Hub
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Ws.WriteMessage(mt, payload)
}

func (c *Connection) Writer() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Ws.Close()
	}()

	// mlog.Info("before write loop")

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				mlog.Warning("Closing socket ...")
				return
			}

			c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Ws.WriteJSON(message); err != nil {
				mlog.Error(err)
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				mlog.Info("error with ping: %s", err.Error())
				return
			}
		}
	}

}

func (c *Connection) Reader() {
	defer func() {
		c.Hub.Unregister <- c
		c.Ws.Close()
	}()

	c.Ws.SetReadLimit(maxMessageSize)
	c.Ws.SetReadDeadline(time.Now().Add(pongWait))
	c.Ws.SetPongHandler(func(string) error {
		c.Ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// mlog.Info("before read loop")

	for {
		_, message, err := c.Ws.ReadMessage()
		if err != nil {
			break
		}

		c.Hub.FromWs <- message
	}

	// for {
	// 	var msgIn dto.MessageIn
	// 	err := c.Ws.ReadJSON(&msgIn)
	// 	if err != nil {
	// 		mlog.Info("error reading socket: %s", err.Error())
	// 		break
	// 	}

	// 	// if DEBUG {
	// 	mlog.Info("client type is: %s", msgIn)
	// 	// }

	// 	//		c.client = msgIn

	// 	msg := &pubsub.Message{}
	// 	c.hub.bus.Pub(msg, msgIn.Topic)

	// 	// switch msgIn.Topic {
	// 	// case "storage:move":
	// 	// 	msg := &pubsub.Message{}
	// 	// 	c.hub.bus.Pub(msg, "cmd.storageMove")

	// 	// case "storage:update":
	// 	// 	msg := &pubsub.Message{}
	// 	// 	c.hub.bus.Pub(msg, "cmd.storageUpdate")

	// 	// default:
	// 	// 	mlog.Info("Unexpected Topic: " + msgIn.Topic)
	// 	// }
	// }

}
