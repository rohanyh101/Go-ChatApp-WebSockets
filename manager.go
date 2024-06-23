package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	webSocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex

	otps RetentionMap

	handlers map[string]EventHandler
}

func NewManager(ctx context.Context) *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
		otps:     NewRetantionMap(ctx, 5*time.Second),
	}

	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessage
}

func SendMessage(event Event, c *Client) error {
	fmt.Println(event)
	return nil
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	}

	return errors.New("there is no sush event type")
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {

	otp := r.URL.Query().Get("otp")
	if otp == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !m.otps.VerifyOTP(otp) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Println("new connection...")

	// update regular HTTP connection to a websocket connection...
	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	client := NewClient(conn, m)
	m.addClient(client)

	// client processes...
	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) loginHandler(w http.ResponseWriter, r *http.Request) {
	type userloginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req userloginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if req.Username == "rohan" && req.Password == "demo" {
		type response struct {
			OTP string `json:"otp"`
		}

		otp := m.otps.NewOTP()
		res := response{
			OTP: otp.Key,
		}

		data, err := json.Marshal(res)
		if err != nil {
			log.Fatal("error marshalling response: ", err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.conn.Close()
		delete(m.clients, client)
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:8080":
		return true

	default:
		return false
	}
}
