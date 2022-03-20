package sse

import (
	"ShoppingList-Backend/internal/pkg/user"
	"ShoppingList-Backend/pkg/middleware"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Client struct {
	User        *user.AppUser
	messageChan chan []byte
}

type Broker struct {
	Notifier       chan *Notification
	newClients     chan *Client
	closingClients chan *Client
	clients        map[*Client]bool
}

type BrokerEvent struct {
	EventType  string
	EventData  any
	Recipients []string
}

type Notification struct {
	Payload []byte
}

func CreateEvent(event BrokerEvent) []byte {
	json, err := json.Marshal(event)
	if err != nil {
		return []byte("")
	}
	return json
}

func NewNotification(payload []byte) *Notification {
	return &Notification{
		Payload: payload,
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

func setupQueues(conn *amqp.Connection, queueName string) error {
	ch, err := conn.Channel()
	if err != nil {
		zap.S().Errorf("Could not create channel: %w", err)
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("sse-events", "fanout", true, false, false, false, nil)
	if err != nil {
		zap.S().Errorf("Could not create exhange: %w", err)
		return err
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		true,      // delete when unused
		false,     // exclusive
		false,     // no wait
		nil,       // arguments
	)
	if err != nil {
		zap.S().Errorf("could not declare queue: %w", err)
		return err
	}

	err = ch.QueueBind(q.Name, "", "sse-events", false, nil)
	if err != nil {
		zap.S().Errorf("Could not bind exchange: %w", err)
		return err
	}

	return nil
}

func (b *Broker) listen(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		zap.S().Errorf("Could not create channel: %w", err)
		return err
	}
	defer ch.Close()
	for {
		select {
		case s := <-b.newClients:
			b.clients[s] = true
			zap.L().Sugar().Infof("Client added. %d registered clients", len(b.clients))
		case s := <-b.closingClients:
			delete(b.clients, s)
			zap.L().Sugar().Infof("Removed client. %d registered clients", len(b.clients))
		case event := <-b.Notifier:
			ch.Publish("sse-events", "", false, false, amqp.Publishing{
				ContentType: "application/json",
				Body:        event.Payload,
			})
		}
	}
}

func (b *Broker) consume(conn *amqp.Connection, queueName string) error {
	ch, err := conn.Channel()
	if err != nil {
		zap.S().Errorf("Could not create channel: %w", err)
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			payload := &BrokerEvent{}
			err = json.Unmarshal(d.Body, payload)
			if err != nil {
				zap.S().Errorf("could not unmarshal rabbitmq message: %w", err)
				continue
			}
			for client := range b.clients {
				for _, userId := range payload.Recipients {
					if client.User.ID == userId {
						client.messageChan <- d.Body
					}
				}
			}
			// zap.S().Infof("Message received: %s", d.Body)
		}
	}()

	<-forever
	return nil
}

func NewBroker(getRabbitMqConn func() (*amqp.Connection, error)) (broker *Broker) {
	broker = &Broker{
		Notifier:       make(chan *Notification, 1),
		newClients:     make(chan *Client),
		closingClients: make(chan *Client),
		clients:        make(map[*Client]bool),
	}
	rabbitMqConn, err := getRabbitMqConn()
	if err != nil {
		zap.S().Errorf("could not get rabbitmq connection :%w", err)
	}

	qId, err := uuid.NewUUID()
	if err != nil {
		zap.S().Errorf("Could not create queue id: %w", err)
		return
	}
	queueName := "sse-events-" + qId.String()

	err = setupQueues(rabbitMqConn, queueName)
	if err != nil {
		zap.S().Errorf("could not setup queues: %w", err)
		return
	}

	go broker.listen(rabbitMqConn)

	go broker.consume(rabbitMqConn, queueName)

	return
}
