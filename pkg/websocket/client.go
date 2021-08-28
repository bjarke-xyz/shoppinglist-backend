package websocket

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
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

var (
	newline = []byte{'\n'}
)

type SessionInfo map[string]interface{}

type Client struct {
	hub *Hub

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	SessionInfo SessionInfo
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	client := &Client{
		hub:         hub,
		conn:        conn,
		send:        make(chan []byte, 256),
		SessionInfo: make(SessionInfo),
	}
	client.hub.register <- client
	return client
}

func (c *Client) ReadWritePump() {
	var wg sync.WaitGroup
	wg.Add(2)
	go c.writePump(&wg)
	go c.readPump(&wg)
	wg.Wait()
}

// readPump pumps messages from the websocket connection to the hub
//
// readPump runs in a per-connection goroutine. This ensures that there is
// at most one reader on a connection by executing all reads from this goroutine
func (c *Client) readPump(wg *sync.WaitGroup) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		wg.Done()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.S().Errorf("Websocket error: %v", err)
			}
			break
		}
		zap.S().Infow("Message received", "client", c.SessionInfo, "size", len(message))
	}
}

// writePump pumps messages from the hub to the websocket connection
//
// a goroutine running writePump is started for each connection.
// There is at most one writer to a connection at a time
func (c *Client) writePump(wg *sync.WaitGroup) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		wg.Done()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Channel was closed by the hub
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				zap.S().Errorf("Could not get websocket writer: %v", err)
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				zap.S().Errorf("Could not close writer: %v", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				zap.S().Errorf("Could not write ping: %v", err)
				return
			}
		}
	}

}
