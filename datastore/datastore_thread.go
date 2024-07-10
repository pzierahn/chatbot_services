package datastore

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/llm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Thread struct {
	// ID of the thread
	Id uuid.UUID `bson:"_id,omitempty"`

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
	coll := service.mongo.Database(DatabaseName).Collection(CollectionThreads)

	if thread.Id == uuid.Nil {
		return errors.New("thread ID is missing")
	}

	filter := bson.M{
		"_id": thread.Id,
	}

	update := bson.M{
		"$set": thread,
	}

	opts := options.Update().SetUpsert(true)
	_, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

// GetThread returns all messages of a thread
func (service *Service) GetThread(ctx context.Context, userId string, threadId uuid.UUID) (*Thread, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionThreads)

	filter := bson.M{
		"_id":     threadId,
		"user_id": userId,
	}

	var thread Thread
	err := coll.FindOne(ctx, filter).Decode(&thread)
	if err != nil {
		return nil, err
	}

	return &thread, nil
}

// GetThreadIDs returns all thread IDs of a collection for a user
func (service *Service) GetThreadIDs(ctx context.Context, userId string, collectionId uuid.UUID) ([]uuid.UUID, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionThreads)

	filter := bson.M{
		"user_id":       userId,
		"collection_id": collectionId,
	}

	opts := &options.FindOptions{
		Projection: bson.M{
			"_id": 1,
		},
		Sort: bson.M{
			"timestamp": -1,
		},
	}

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	defer func() { _ = cursor.Close(ctx) }()

	var result []uuid.UUID
	for cursor.Next(ctx) {
		var thread Thread
		err = cursor.Decode(&thread)
		if err != nil {
			return nil, err
		}

		result = append(result, thread.Id)
	}

	return result, nil
}
