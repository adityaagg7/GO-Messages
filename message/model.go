package message

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Body     string             `bson:"body,omitempty" json:"body"`
	RoomID   string             `bson:"room_id,omitempty" json:"room_id"`
	SenderID string             `bson:"sender_id,omitempty" json:"sender_id"`
}
