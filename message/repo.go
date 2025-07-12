package message

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type MessageRepo interface {
	PostMessage(ctx context.Context, msg *Message) (*Message, error)
	GetMessagesByRoomId(ctx context.Context, roomID primitive.ObjectID) ([]Message, error)
}

type MessageRepoImpl struct {
	messageCollection *mongo.Collection
}

func NewMessageRepository(client *mongo.Client) MessageRepo {
	return &MessageRepoImpl{
		messageCollection: client.Database(os.Getenv("MONGO_DB_NAME")).Collection("messages"),
	}
}

func (r *MessageRepoImpl) PostMessage(ctx context.Context, msg *Message) (*Message, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	log.Println("Post Message: ", msg)
	result, err := r.messageCollection.InsertOne(timeoutCtx, msg)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		msg.ID = oid
	}
	return msg, nil
}

// GetMessagesByRoomId retrieves all messages for a given room ID from the database.
// Returns the list of messages or an error if any issue occurs during the operation.
func (r *MessageRepoImpl) GetMessagesByRoomId(ctx context.Context, roomID primitive.ObjectID) ([]Message, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"room_id": roomID.Hex()}

	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}})
	cursor, err := r.messageCollection.Find(timeoutCtx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Default().Println(err.Error())
		}
	}(cursor, timeoutCtx)

	var messages []Message
	if err := cursor.All(timeoutCtx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
