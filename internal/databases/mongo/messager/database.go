package messager

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

// Client is a global variable that holds the MongoDB client instance for database connections.
var Client *mongo.Client

// ConnectDB initializes a connection to the MongoDB database and returns a reference to the MongoDB client instance.
func ConnectDB() *mongo.Client {
	if Client != nil {
		return Client
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB not responding:", err)
	}

	log.Println("Connected to MongoDB")
	Client = client
	return client
}

// DB is a global variable initialized with a MongoDB client connection using the ConnectDB function.
var DB = ConnectDB()

// GetCollection retrieves a MongoDB collection by its name from the database specified in the MONGO_DB_NAME environment variable.
func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database(os.Getenv("MONGO_DB_NAME")).Collection(collectionName)
}
