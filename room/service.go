package room

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"messages-go/models/request"
	"messages-go/utils"
	"strings"
)

var (
	ErrRoomNotFound     = errors.New("room not found")
	ErrMongoWriteFailed = errors.New("mongo write failed")
)

// RoomService defines the interface for managing room operations, including creation and retrieval of rooms.
type RoomService interface {
	CreateRoom(ctx context.Context, req request.CreateRoomRequest) (*Room, error)
	GetRoom(ctx context.Context, id string) (*Room, error)
	UpdateRoomName(ctx context.Context, id string, name string) (*Room, error)
}

// RoomServiceImpl is a service that handles business logic related to room operations using a room repository.
type RoomServiceImpl struct {
	roomRepo RepoRoom
}

// NewRoomService initializes and returns a new instance of RoomServiceImpl with the provided room repository.
func NewRoomService(roomRepo RepoRoom) *RoomServiceImpl {
	return &RoomServiceImpl{roomRepo: roomRepo}
}

// CreateRoom handles the creation of a new room, automatically generating a name if none is provided in the request.
func (rs *RoomServiceImpl) CreateRoom(ctx context.Context, req request.CreateRoomRequest) (*Room, error) {
	var roomName string
	if req.Name == nil || strings.TrimSpace(*req.Name) == "" {
		roomName = utils.GenerateRoomName()
	} else {
		roomName = *req.Name
	}

	room := Room{Name: roomName}
	return rs.roomRepo.CreateRoom(ctx, &room)
}

// GetRoom retrieves a room by its unique identifier from the repository and returns the room or an error if not found.
func (rs *RoomServiceImpl) GetRoom(ctx context.Context, id string) (*Room, error) {
	room, err := rs.roomRepo.GetRoomByID(ctx, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrRoomNotFound
	}
	return room, err
}

// UpdateRoomName updates the name of an existing room by its ID in the repository and returns the updated room or an error.
func (rs *RoomServiceImpl) UpdateRoomName(ctx context.Context, id string, name string) (*Room, error) {
	updatedRoom, err := rs.roomRepo.UpdateRoomName(ctx, id, name)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrRoomNotFound
		}
		return nil, ErrMongoWriteFailed
	}
	return updatedRoom, nil

}
