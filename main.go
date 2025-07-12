package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"messages-go/internal/databases/mongo/messager"
	"messages-go/routes"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// main is the entry point of the application, initializing the database, setting up routes, and starting the HTTP server.
func main() {
	// Initialize database
	var db = messager.ConnectDB()

	// Create a new Fiber instance
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Content-Type",
	}))
	// Set up routes
	routes.SetupRoutes(app, db)

	// Create a channel to listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Run the server in a goroutine
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for termination signal
	msg := <-quit
	log.Printf("Shutting down server due to %s\n", msg)

	// Create a context with timeout to allow cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shutdown the Fiber app
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
