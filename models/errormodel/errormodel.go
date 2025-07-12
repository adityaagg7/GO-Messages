package errormodel

import "errors"

var (
	ErrRoomNotFound     = errors.New("room not found")
	ErrMessagesNotFound = errors.New("no messages found")
	ErrMongoWriteFailed = errors.New("mongo write failed")
)
