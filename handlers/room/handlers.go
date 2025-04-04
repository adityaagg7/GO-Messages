package room

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"messages-go/models/request"
	"messages-go/models/response"
	"messages-go/services/room"
)

// CreateRoom handles HTTP requests to create a new room, parsing the request body and returning a response with the result.
func CreateRoom(c *fiber.Ctx) error {

	var req request.CreateRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  400,
			Message: "Invalid Request Body",
		})
	}
	log.Info("Create Room Request Received.")

	roomResp, err := room.CreateRoom(req)
	if err != nil {
		return c.Status(500).JSON(response.APIResponse{
			Error:   err.Error(),
			Status:  500,
			Message: "Failed To Create Room",
		})
	}

	return c.Status(201).JSON(response.APIResponse{
		Status:  201,
		Message: "Room Created",
		Data:    roomResp,
	})
}
