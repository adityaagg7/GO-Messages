package room

import (
	room_model "messages-go/models/mongo_messager/room"
	"messages-go/models/request"
	room_repo "messages-go/repos/mongo_messager/room"
	"messages-go/utils"
	"strings"
)

// CreateRoom creates a new room with a specified name or generates a random name if none is provided.
// It returns the created room's ID or an error if the operation fails.
func CreateRoom(req request.CreateRoomRequest) (*room_model.Room, error) {
	var roomName string
	if req.Name == nil || strings.TrimSpace(*req.Name) == "" {
		roomName = utils.GenerateRoomName()
	} else {
		roomName = *req.Name
	}

	room := room_model.Room{Name: roomName}
	return room_repo.CreateRoom(&room)
}
