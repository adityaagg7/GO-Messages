package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub maintains the set of active connections and broadcasts messages to the connections.
type Hub struct {
	// Registered connections mapped by room ID
	rooms map[string]map[*Client]bool

	// Inbound messages from the connections.
	broadcast chan BroadcastMessage

	// Register requests from the connections.
	register chan *Client

	// Unregister requests from connections.
	unregister chan *Client

	// Mutex to protect the rooms map
	mu sync.RWMutex
}

type BroadcastMessage struct {
	RoomID  string      `json:"room_id"`
	Message interface{} `json:"message"`
}

// Global hub instance
var GlobalHub = &Hub{
	rooms:      make(map[string]map[*Client]bool),
	broadcast:  make(chan BroadcastMessage),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

// Start initializes and runs the hub
func (h *Hub) Start() {
	go h.run()
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			h.mu.Unlock()
			log.Printf("Client registered for room: %s", client.RoomID)

		case client := <-h.unregister:
			h.mu.Lock()
			if room, ok := h.rooms[client.RoomID]; ok {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.Send)
					if len(room) == 0 {
						delete(h.rooms, client.RoomID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client unregistered for room: %s", client.RoomID)

		case message := <-h.broadcast:
			h.mu.RLock()
			room := h.rooms[message.RoomID]
			h.mu.RUnlock()

			if room != nil {
				messageBytes, err := json.Marshal(message.Message)
				if err != nil {
					log.Printf("Error marshaling message: %v", err)
					continue
				}

				for client := range room {
					select {
					case client.Send <- messageBytes:
					default:
						close(client.Send)
						delete(room, client)
					}
				}
			}

		}
	}
}

// BroadcastToRoom sends a message to all connections in a specific room
func (h *Hub) BroadcastToRoom(roomID string, message interface{}) {
	h.broadcast <- BroadcastMessage{
		RoomID:  roomID,
		Message: message,
	}
}

// GetRoomConnections returns the number of active connections in a room
func (h *Hub) GetRoomConnections(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, ok := h.rooms[roomID]; ok {
		return len(room)
	}
	return 0
}
