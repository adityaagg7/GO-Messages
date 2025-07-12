package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"messages-go/message"
	"messages-go/room"
	ws "messages-go/websocket"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(app *fiber.App, client *mongo.Client) {
	// Initialize WebSocket hub
	hub := ws.GlobalHub
	hub.Start()
	wsHandler := ws.NewHandler(hub)

	// Initialize REST handlers
	roomHandler, roomRepo, _ := room.InitRoomHandler(client)
	messageHandler, _, _ := message.InitMessageHandler(client, roomRepo, wsHandler)

	// Pass WebSocket handler to message handler for broadcasting
	// You'll need to modify your message handler to accept this

	// API routes
	api := app.Group("/api")
	setupRoomRoutes(api, roomHandler)
	setupMessageRoutes(api, messageHandler)

	// WebSocket routes
	setupWebSocketRoutes(app, wsHandler)
}

func setupRoomRoutes(api fiber.Router, handler room.RoomHandler) {
	roomGroup := api.Group("/room")
	roomGroup.Post("/", handler.CreateRoom)
	roomGroup.Get("/:name", handler.GetRoom)
	roomGroup.Patch("/:id", handler.UpdateRoomName)
}

func setupMessageRoutes(api fiber.Router, handler message.MessageHandler) {
	messageGroup := api.Group("/message")
	messageGroup.Post("/", handler.PostMessage)
	messageGroup.Get("/:roomId", handler.GetMessages)
}

func setupWebSocketRoutes(app *fiber.App, wsHandler *ws.Handler) {
	// WebSocket upgrade middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint
	app.Get("/ws/:roomId", websocket.New(wsHandler.HandleConnection))
}
