package datastore

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/llm"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Thread struct {
	// ID of the message
	Id uuid.UUID `bson:"_id,omitempty"`

	// Thread ID
	ThreadId uuid.UUID `bson:"thread_id,omitempty"`

	// User ID
	UserId string `bson:"user_id,omitempty"`

	// Collection ID
	CollectionId uuid.UUID `bson:"collection_id,omitempty"`

	// Timestamp of the last message
	Timestamp time.Time `bson:"timestamp,omitempty"`

	// Messages
	Messages []*llm.Message `bson:"messages,omitempty"`
}

// StoreThread stores a thread
func (service *Service) StoreThread(ctx context.Context, thread *Thread) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionMessages)

	_, err := coll.InsertOne(ctx, thread)
	if err != nil {
		return err
	}

	return nil
}

// GetThread returns all messages of a thread
func (service *Service) GetThread(ctx context.Context, userId string, threadId uuid.UUID) (*Thread, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionMessages)

	filter := bson.M{
		"thread_id": threadId,
		"user_id":   userId,
	}

	var thread Thread
	err := coll.FindOne(ctx, filter).Decode(&thread)
	if err != nil {
		return nil, err
	}

	return &thread, nil
}
