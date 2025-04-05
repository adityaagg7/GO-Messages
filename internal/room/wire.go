package room

import (
	"go.mongodb.org/mongo-driver/mongo"
	roomHandler "messages-go/handlers/room"
	roomRepo "messages-go/repos/mongo_messager/room"
	roomService "messages-go/services/room"
)

func InitRoomHandler(client *mongo.Client) roomHandler.RoomHandler {
	repo := roomRepo.NewRoomRepository(client)
	service := roomService.NewRoomService(repo)
	handler := roomHandler.NewRoomHandler(service)
	return handler
}
