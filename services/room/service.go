package room

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	room_model "messages-go/models/mongo_messager/room"
	"messages-go/models/request"
	room_repo "messages-go/repos/mongo_messager/room"
	"messages-go/utils"
	"strings"
)

var ErrRoomNotFound = errors.New("room not found")

// RoomService defines the interface for managing room operations, including creation and retrieval of rooms.
type RoomService interface {
	CreateRoom(ctx context.Context, req request.CreateRoomRequest) (*room_model.Room, error)
	GetRoom(ctx context.Context, id string) (*room_model.Room, error)
}

// RoomServiceImpl is a service that handles business logic related to room operations using a room repository.
type RoomServiceImpl struct {
	roomRepo room_repo.RepoRoom
}

// NewRoomService initializes and returns a new instance of RoomServiceImpl with the provided room repository.
func NewRoomService(roomRepo room_repo.RepoRoom) *RoomServiceImpl {
	return &RoomServiceImpl{roomRepo: roomRepo}
}

// CreateRoom handles the creation of a new room, automatically generating a name if none is provided in the request.
func (rs *RoomServiceImpl) CreateRoom(ctx context.Context, req request.CreateRoomRequest) (*room_model.Room, error) {
	var roomName string
	if req.Name == nil || strings.TrimSpace(*req.Name) == "" {
		roomName = utils.GenerateRoomName()
	} else {
		roomName = *req.Name
	}

	room := room_model.Room{Name: roomName}
	return rs.roomRepo.CreateRoom(ctx, &room)
}

// GetRoom retrieves a room by its unique identifier from the repository and returns the room or an error if not found.
func (rs *RoomServiceImpl) GetRoom(ctx context.Context, id string) (*room_model.Room, error) {
	room, err := rs.roomRepo.GetRoomByID(ctx, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrRoomNotFound
	}
	return room, err
}
