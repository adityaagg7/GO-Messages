package websocket

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

// Client represents a WebSocket connection for a specific room
type Client struct {
	Conn   *websocket.Conn
	RoomID string
	Send   chan []byte
	Hub    *Hub
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, roomID string, hub *Hub) *Client {
	return &Client{
		Conn:   conn,
		RoomID: roomID,
		Send:   make(chan []byte, 256),
		Hub:    hub,
	}
}

// Start begins the client's read and write pumps
func (c *Client) Start() {
	log.Printf("Starting client for room: %s", c.RoomID)
	c.Hub.register <- c
	log.Printf("Client registered, starting pumps for room: %s", c.RoomID)

	go c.writePump()
	log.Printf("Write pump started for room: %s", c.RoomID)

	c.readPump() // This should block here
	log.Printf("Read pump exited for room: %s", c.RoomID)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		// We're only broadcasting, not receiving messages through WebSocket
		// Messages are still posted through the REST API
	}
}
