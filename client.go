package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait = 10 * time.Second

	// in this case we care waiting for 90% of the pongWait time before sending a ping again...
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[*client]bool

type client struct {
	connection *websocket.Conn
	manager    *Manager

	// egress is used to avoid concurrent writes to the websocket connection
	egress chan Event
}

func NewClient(conn *websocket.Conn, m *Manager) *client {
	return &client{
		connection: conn,
		manager:    m,
		egress:     make(chan Event),
	}
}

func (c *client) readMessages() {
	defer func() {
		// clean up client connection
		c.manager.removeClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("error setting read deadline: %v\n", err)
		return
	}

	// this depends on your business logic...
	// this is known as "jumbo frames" where the maximum size of the message is 512 bytes
	c.connection.SetReadLimit(512)

	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v\n", err)
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error unmarshalling message: %v\n", err)
			continue
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Printf("error routing event: %v\n", err)
		}
	}
}

func (c *client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("error writing close message: %v\n", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("error marshalling message: %v\n", err)
				continue
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("error writing message: %v\n", err)
			}
			log.Println("message sent to client")

		case <-ticker.C:
			log.Println("ping")

			// sending ping message to client
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Printf("error writing ping message: %v\n", err)
				return
			}
		}
	}
}

// don't forgot to reset the read deadline after receiving a pong message...
// otherwise the connection will be closed after pongWait time
func (c *client) pongHandler(pongMsg string) error {
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
