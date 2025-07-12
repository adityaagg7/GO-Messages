package message

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
	"messages-go/models/errormodel"
	"messages-go/models/response"
	ws "messages-go/websocket"
	"strings"
)

type MessageHandler interface {
	PostMessage(c *fiber.Ctx) error
	GetMessages(c *fiber.Ctx) error
}

type MessageHandlerImpl struct {
	messageService MessageService
	wsHandler      *ws.Handler
}

func NewMessageHandler(messageService MessageService, webSocketHandler *ws.Handler) MessageHandler {
	return &MessageHandlerImpl{messageService: messageService, wsHandler: webSocketHandler}
}

func (mh *MessageHandlerImpl) PostMessage(c *fiber.Ctx) error {
	var postMessageRequest Message
	if err := c.BodyParser(&postMessageRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Request Body.",
		})
	}
	log.Println("Post Message Request Received:", postMessageRequest)
	message, err := mh.messageService.PostMessage(c.Context(), &postMessageRequest)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  fiber.StatusInternalServerError,
			Message: "Failed To Post Message.",
		})
	}

	if mh.wsHandler != nil {
		mh.wsHandler.BroadcastToRoom(message.RoomID, map[string]interface{}{
			"type":    "new_message",
			"message": message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.APIResponse{
		Data:    message,
		Status:  fiber.StatusOK,
		Message: "Message Posted.",
	})
}

func (mh *MessageHandlerImpl) GetMessages(c *fiber.Ctx) error {

	var roomName = c.Params("roomName")
	if strings.TrimSpace(roomName) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.APIResponse{
			Error:   "Missing room name in path",
			Status:  fiber.StatusBadRequest,
			Message: "roomName Path Variable is required to be not empty.",
		})
	}

	log.Println("Get Messages from Room with name: ", roomName, " Request Received.")
	getMessageResp, err := mh.messageService.GetMessages(c.Context(), roomName)

	if errors.Is(err, errormodel.ErrMessagesNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  fiber.StatusNotFound,
			Message: "No messages Found with given roomname.",
		})
	} else if errors.Is(err, errormodel.ErrRoomNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  fiber.StatusNotFound,
			Message: "No Room found given roomname.",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  fiber.StatusInternalServerError,
			Message: "Failed To Get Messages",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Messages Found",
		Data:    getMessageResp,
	})
}
