package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*client]bool

type client struct {
	connection *websocket.Conn
	manager    *Manager

	// egress is used to avoid concurrent writes to the websocket connection
	egress chan []byte
}

func NewClient(conn *websocket.Conn, m *Manager) *client {
	return &client{
		connection: conn,
		manager:    m,
		egress:     make(chan []byte),
	}
}

func (c *client) readMessages() {
	defer func() {
		// clean up client connection
		c.manager.removeClient(c)
	}()

	for {
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v\n", err)
			}
			break
		}

		for wsclient := range c.manager.clients {
			wsclient.egress <- payload
		}

		log.Println("message type: ", messageType)
		log.Println("message: ", string(payload))
	}
}

func (c *client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("error writing close message: %v\n", err)
				}
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("error writing message: %v\n", err)
			}
			log.Println("message sent: ", string(message))
		}
	}
}
