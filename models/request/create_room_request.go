package request

// CreateRoomRequest represents the structure for requests to create a new room, optionally specifying a room name.
type CreateRoomRequest struct {
	Name *string `json:"name"`
}
