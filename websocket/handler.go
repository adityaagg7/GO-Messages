package websocket

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

// Handler manages WebSocket connections
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

// HandleConnection handles WebSocket connections
func (h *Handler) HandleConnection(c *websocket.Conn) {
	roomID := c.Params("roomId")
	if roomID == "" {
		log.Println("Room ID not provided")
		c.Close()
		return
	}

	// Create and start client
	client := NewClient(c, roomID, h.hub)
	client.Start()
}

// BroadcastToRoom is a convenience method to broadcast to a specific room
func (h *Handler) BroadcastToRoom(roomID string, message interface{}) {
	h.hub.BroadcastToRoom(roomID, message)
}

// GetRoomConnections returns the number of active connections in a room
func (h *Handler) GetRoomConnections(roomID string) int {
	return h.hub.GetRoomConnections(roomID)
}
