package sse

import (
	"ShoppingList-Backend/internal/pkg/user"
	"ShoppingList-Backend/pkg/middleware"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Client struct {
	User        *user.AppUser
	messageChan chan []byte
}
type ClientFilter func(*Client) bool

type Broker struct {
	Notifier       chan *Notification
	newClients     chan *Client
	closingClients chan *Client
	clients        map[*Client]bool
}

type BrokerEvent[T any] struct {
	EventType string
	EventData T
}

type Notification struct {
	Payload []byte
	Filter  ClientFilter
}

func CreateEvent[T any](event BrokerEvent[T]) []byte {
	json, err := json.Marshal(event)
	if err != nil {
		return []byte("")
	}
	return json
}

func NewNotification(payload []byte, filter ClientFilter) *Notification {
	return &Notification{
		Payload: payload,
		Filter:  filter,
	}
}

func (b *Broker) Handle(w http.ResponseWriter, req *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	client := &Client{
		User:        middleware.UserFromContext(req.Context()),
		messageChan: make(chan []byte),
	}
	b.newClients <- client
	defer func() {
		b.closingClients <- client
	}()

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		b.closingClients <- client
	}()

	for {
		fmt.Fprintf(w, "data: %s\n\n", <-client.messageChan)
		flusher.Flush()
	}
}

func (b *Broker) listen() {
	for {
		select {
		case s := <-b.newClients:
			b.clients[s] = true
			zap.L().Sugar().Infof("Client added. %d registered clients", len(b.clients))
		case s := <-b.closingClients:
			delete(b.clients, s)
			zap.L().Sugar().Infof("Removed client. %d registered clients", len(b.clients))
		case event := <-b.Notifier:
			for client := range b.clients {
				if event.Filter(client) {
					client.messageChan <- event.Payload
				}
			}
		}
	}
}

func NewBroker() (broker *Broker) {
	broker = &Broker{
		Notifier:       make(chan *Notification, 1),
		newClients:     make(chan *Client),
		closingClients: make(chan *Client),
		clients:        make(map[*Client]bool),
	}
	go broker.listen()
	return
}
