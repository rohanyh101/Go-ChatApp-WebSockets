package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket connection established")

	// upgrade regular HTTP connection to a websocket connection
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, m)

	m.addClient(client)

	// strat 2 go routines to read and write messages
	// read messages from the client
	go client.readMessages()

	// write messages to the client
	go client.writeMessages()
}

func (m *Manager) addClient(c *client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[c]; !ok {
		m.clients[c] = true
	}
}

func (m *Manager) removeClient(c *client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[c]; ok {
		c.connection.Close()
		delete(m.clients, c)
	}
}
