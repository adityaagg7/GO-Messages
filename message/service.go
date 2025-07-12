package message

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"messages-go/models/errormodel"
	"messages-go/room"
)

type MessageService interface {
	PostMessage(ctx context.Context, msg *Message) (*Message, error)
	GetMessages(ctx context.Context, roomId string) ([]Message, error)
}

type MessageServiceImpl struct {
	messageRepo MessageRepo
	roomRepo    room.RoomRepo
}

func NewMessageService(messageRepo MessageRepo, roomRepo room.RoomRepo) *MessageServiceImpl {
	return &MessageServiceImpl{
		messageRepo: messageRepo,
		roomRepo:    roomRepo,
	}
}

func (ms *MessageServiceImpl) PostMessage(ctx context.Context, msg *Message) (*Message, error) {
	_, err := ms.roomRepo.GetRoomByID(ctx, msg.RoomID)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errormodel.ErrRoomNotFound
	}
	log.Println("Posting Message: ", msg)
	return ms.messageRepo.PostMessage(ctx, msg)
}

func (ms *MessageServiceImpl) GetMessages(ctx context.Context, roomId string) ([]Message, error) {
	roomData, err := ms.roomRepo.GetRoomByID(ctx, roomId)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errormodel.ErrRoomNotFound
	}
	log.Println("Getting Messages for: ", roomData.ID.String())
	messageList, err := ms.messageRepo.GetMessagesByRoomId(ctx, roomData.ID)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errormodel.ErrMessagesNotFound
	}
	return messageList, err
}
