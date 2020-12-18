package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type client struct {
	conn *websocket.Conn
	send chan []byte
	hub  *hub
	file *file
}

func (c *client) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.hub.broadcast <- msg
	}
}

func (c *client) write() {
	defer c.conn.Close()

	for {
		select {
		case msg, _ := <-c.send:
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				break
			}

			c.file.write(msg)
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(h *hub, f *file, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &client{
		conn: conn,
		send: make(chan []byte, 256),
		hub:  h,
		file: f,
	}
	client.hub.register <- client

	go client.write()
	go client.read()
}
