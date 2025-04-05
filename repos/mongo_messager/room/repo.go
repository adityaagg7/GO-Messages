package room

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"messages-go/models/mongo_messager/room"
	"os"
	"time"
)

// RepoRoom defines an interface for room persistence operations including create, retrieve, and update functionalities.
type RepoRoom interface {
	CreateRoom(ctx context.Context, room *room.Room) (*room.Room, error)
	GetRoomByID(ctx context.Context, id string) (*room.Room, error)
	UpdateRoomName(ctx context.Context, id string, name string) error
}

// RepoRoomImpl is a concrete implementation of the RepoRoom interface.
// It interacts with the MongoDB collection to manage room data.
type RepoRoomImpl struct {
	roomCollection *mongo.Collection
}

// NewRoomRepository initializes and returns a new instance of RepoRoom for managing room data in MongoDB.
func NewRoomRepository(client *mongo.Client) RepoRoom {
	return &RepoRoomImpl{
		roomCollection: client.Database(os.Getenv("MONGO_DB_NAME")).Collection("rooms"),
	}
}

// CreateRoom inserts a new room document into the database and returns the created room or an error if the operation fails.
func (r *RepoRoomImpl) CreateRoom(ctx context.Context, rm *room.Room) (*room.Room, error) {
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
func (r *RepoRoomImpl) GetRoomByID(ctx context.Context, id string) (*room.Room, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var rm room.Room
	err = r.roomCollection.FindOne(timeoutCtx, bson.M{"_id": objID}).Decode(&rm)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	return &rm, err
}

// UpdateRoomName updates the name of a room in the database based on the provided room ID and new name.
func (r *RepoRoomImpl) UpdateRoomName(ctx context.Context, id string, name string) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	_, err = r.roomCollection.UpdateOne(
		timeoutCtx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"name": name}},
	)

	return err
}
