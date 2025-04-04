package routes

import (
	"github.com/gofiber/fiber/v2"
	"messages-go/handlers/room"
)

// SetupRoutes configures the routes for the application, grouping them under "/api" and defining endpoints for rooms.
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	roomGroup := api.Group("/room")

	roomGroup.Post("/", room.CreateRoom)
}
