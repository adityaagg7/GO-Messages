package room

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
	"messages-go/models/errormodel"
	"messages-go/models/request"
	"messages-go/models/response"
	"strings"
)

// RoomHandler defines the interface for handling HTTP requests related to room operations, such as creation and retrieval.
type RoomHandler interface {
	CreateRoom(c *fiber.Ctx) error
	GetRoom(c *fiber.Ctx) error
	UpdateRoomName(c *fiber.Ctx) error
}

// RoomHandlerImpl implements the RoomHandler interface and handles HTTP requests related to room operations.
type RoomHandlerImpl struct {
	roomService RoomService
}

// NewRoomHandler initializes and returns a new RoomHandler with the provided RoomService implementation.
func NewRoomHandler(roomService RoomService) RoomHandler {
	return &RoomHandlerImpl{roomService: roomService}
}

// CreateRoom handles the creation of a new room by parsing the request body, invoking the service layer, and returning a response.
func (rh *RoomHandlerImpl) CreateRoom(c *fiber.Ctx) error {
	var req request.CreateRoomRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Request Body",
		})
	}

	log.Println("Create Room Request Received.")

	roomResp, err := rh.roomService.CreateRoom(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  fiber.StatusInternalServerError,
			Message: "Failed To Create Room",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.APIResponse{
		Status:  fiber.StatusCreated,
		Message: "Room Created",
		Data:    roomResp,
	})
}

// GetRoom handles the retrieval of a room by its ID from the path parameter and returns the room information as a response.
func (rh *RoomHandlerImpl) GetRoom(c *fiber.Ctx) error {
	var roomName string

	roomName = c.Params("name")

	if strings.TrimSpace(roomName) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.APIResponse{
			Error:   "Missing room name in path",
			Status:  fiber.StatusBadRequest,
			Message: "name Path Variable is required to be not empty.",
		})
	}

	log.Println("Get Room with name: ", roomName, " Request Received.")

	getRoomResp, err := rh.roomService.GetRoom(c.Context(), roomName)

	if errors.Is(err, errormodel.ErrRoomNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  fiber.StatusNotFound,
			Message: "No Room Found with given name.",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  fiber.StatusInternalServerError,
			Message: "Failed To Get Room",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Room Found",
		Data:    getRoomResp,
	})
}

// UpdateRoomName handles updating the name of an existing room using the room ID and the new name provided in the request body.
func (rh *RoomHandlerImpl) UpdateRoomName(c *fiber.Ctx) error {
	roomId := c.Params("id")
	var req request.UpdateRoomRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Request Body",
		})
	}

	log.Println("Update Room with id ", roomId, " Request Received.")

	roomResp, err := rh.roomService.UpdateRoomName(c.Context(), roomId, *req.Name)

	if errors.Is(err, errormodel.ErrRoomNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  fiber.StatusNotFound,
			Message: "No Room Found with given id.",
		})
	} else if errors.Is(err, errormodel.ErrMongoWriteFailed) {
		return c.Status(fiber.StatusInternalServerError).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  fiber.StatusInternalServerError,
			Message: "Failed To Write Updated Room",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  fiber.StatusInternalServerError,
			Message: "Failed To Update Room",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Room Updated",
		Data:    roomResp,
	})

}
