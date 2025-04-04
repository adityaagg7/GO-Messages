package room

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"messages-go/databases/mongo/messager"
	"messages-go/models/mongo_messager/room"
	"time"
)

// roomCollection references the MongoDB collection "room" used to store and manage data related to rooms.
var roomCollection = messager.GetCollection("room")

// ctx is the base context used as a parent for creating derived contexts with specific timeouts or cancellation.
var ctx = context.Background()

// CreateRoom inserts a new room into the room collection and returns the generated room ID or an error if it fails.
func CreateRoom(room *room.Room) (*room.Room, error) {
	var timeOutContext, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result, err := roomCollection.InsertOne(timeOutContext, room)
	if err != nil {
		return nil, err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		room.ID = oid
	}
	return room, nil
}

// GetRoomByID retrieves a room document from the database by its ID and returns the corresponding Room object or an error.
func GetRoomByID(id string) (*room.Room, error) {

	var timeOutContext, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var room room.Room
	err = roomCollection.FindOne(timeOutContext, bson.M{"_id": objID}).Decode(&room)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &room, err
}

// UpdateRoomName updates the name of a room identified by the given ID in the database and returns an error if any occurs.
func UpdateRoomName(id string, name string) error {
	var timeOutContext, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}
	_, err = roomCollection.UpdateOne(timeOutContext, bson.M{"_id": objID}, bson.M{"$set": bson.M{"name": name}})
	return err
}
