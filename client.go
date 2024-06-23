package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	conn    *websocket.Conn
	manager *Manager

	// egress is used to avoid concurrent writes to the websocket connection...
	egress chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		conn:    conn,
		manager: manager,
		egress:  make(chan Event),
	}
}

func (c *Client) readMessages() {
	defer func() {
		// clean up client connection...
		c.manager.removeClient(c)
	}()

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println("error setting read deadline: ", err)
		return
	}

	// set read limit to avoid abuse of the connection...
	c.conn.SetReadLimit(512)
	c.conn.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			log.Println("error unmarshalling message: ", err)
			break
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("error routing event: ", err)
			break
		}
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	tiker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println("error marshalling message: ", err)
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("error writing message: ", err)
				return
			}

			log.Println("message sent...")

		case <-tiker.C:
			log.Println("ping...")

			// send ping message to client...
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte("")); err != nil {
				log.Println("error sending ping message: ", err)
				return
			}
		}
	}
}

func (c *Client) pongHandler(pingMsg string) error {
	log.Println("pong...")
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}
