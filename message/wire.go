package message

import (
	"go.mongodb.org/mongo-driver/mongo"
	"messages-go/room"
	ws "messages-go/websocket"
)

func InitMessageHandler(client *mongo.Client, roomRepo room.RoomRepo, wsHandler *ws.Handler) (MessageHandler, MessageRepo, MessageService) {
	repo := NewMessageRepository(client)
	service := NewMessageService(repo, roomRepo)
	handler := NewMessageHandler(service, wsHandler)
	return handler, repo, service
}
