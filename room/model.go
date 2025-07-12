package room

import "go.mongodb.org/mongo-driver/bson/primitive"

// Room represents a struct containing information about a room, including its ID and name.
type Room struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name,omitempty," json:"name"`
}
