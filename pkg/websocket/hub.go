package websocket

import (
	"encoding/json"

	"go.uber.org/zap"
)

type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Register requests from the clients
	register chan *Client

	// Unregister requests from the clients
	unregister chan *Client
}

type EventPayload struct {
	EventType string      `json:"eventType"`
	EventData interface{} `json:"eventData"`
}
type ClientFilter func(SessionInfo) bool

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		}
	}
}

func (h *Hub) EmitEvent(payload EventPayload, clientFilter ClientFilter) {
	for client := range h.clients {
		if !clientFilter(client.SessionInfo) {
			return
		}
		payloadJson, err := json.Marshal(payload)
		if err != nil {
			zap.S().Errorf("Error marshaling event payload: %v", err)
		}
		client.send <- payloadJson
	}
}
