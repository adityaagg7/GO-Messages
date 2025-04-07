package request

// UpdateRoomRequest represents a request to update the details of a room, primarily its name.
type UpdateRoomRequest struct {
	Name *string `json:"name"`
}
