package room

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func InitRoomHandler(client *mongo.Client) (RoomHandler, RoomRepo, RoomService) {
	repo := NewRoomRepository(client)
	service := NewRoomService(repo)
	handler := NewRoomHandler(service)
	return handler, repo, service
}
