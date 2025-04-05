package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"messages-go/internal/room"
)

// SetupRoutes configures the routes for the application, grouping them under "/api" and defining endpoints for rooms.
func SetupRoutes(app *fiber.App, client *mongo.Client) {
	api := app.Group("/api")
	roomGroup := api.Group("/room")

	var roomHandler = room.InitRoomHandler(client)

	roomGroup.Post("/", roomHandler.CreateRoom)
	roomGroup.Get("/:id", roomHandler.GetRoom)
}
