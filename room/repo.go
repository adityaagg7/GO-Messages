package room

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

// RoomRepo defines an interface for room persistence operations including create, retrieve, and update functionalities.
type RoomRepo interface {
	CreateRoom(ctx context.Context, room *Room) (*Room, error)
	GetRoomByID(ctx context.Context, id string) (*Room, error)
	UpdateRoomName(ctx context.Context, id string, name string) (*Room, error)
}

// RoomRepoImpl is a concrete implementation of the RoomRepo interface.
// It interacts with the MongoDB collection to manage room data.
type RoomRepoImpl struct {
	roomCollection *mongo.Collection
}

// NewRoomRepository initializes and returns a new instance of RoomRepo for managing room data in MongoDB.
func NewRoomRepository(client *mongo.Client) RoomRepo {
	return &RoomRepoImpl{
		roomCollection: client.Database(os.Getenv("MONGO_DB_NAME")).Collection("rooms"),
	}
}

// CreateRoom inserts a new room document into the database and returns the created room or an error if the operation fails.
func (r *RoomRepoImpl) CreateRoom(ctx context.Context, rm *Room) (*Room, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.roomCollection.InsertOne(timeoutCtx, rm)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		rm.ID = oid
	}
	return rm, nil
}

// GetRoomByID retrieves a room by its ID from the database.
// Returns the room or nil if not found, and an error if any issue occurs during the operation.
func (r *RoomRepoImpl) GetRoomByID(ctx context.Context, id string) (*Room, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var rm Room
	err = r.roomCollection.FindOne(timeoutCtx, bson.M{"_id": objID}).Decode(&rm)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	return &rm, err
}

// UpdateRoomName updates the name of a room identified by its ID in the database and returns the updated room or an error.
func (r *RoomRepoImpl) UpdateRoomName(ctx context.Context, id string, name string) (*Room, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var updatedRoom Room

	err = r.roomCollection.FindOneAndUpdate(
		timeoutCtx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"name": name}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedRoom)

	if err != nil {
		return nil, err
	}
	return &updatedRoom, nil
}
