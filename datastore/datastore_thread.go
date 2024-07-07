package datastore

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// AddMessage adds a message to the database
func (service *Service) AddMessage(ctx context.Context, message *Message) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionMessages)

	_, err := coll.InsertOne(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

// AddMessages adds multiple messages to the database
func (service *Service) AddMessages(ctx context.Context, messages []*Message) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionMessages)

	var docs []interface{}
	for _, message := range messages {
		docs = append(docs, message)
	}

	_, err := coll.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	return nil
}

// GetMessages gets all messages from a thread
func (service *Service) GetMessages(ctx context.Context, userId string, threadId uuid.UUID) ([]*Message, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionMessages)

	filter := bson.M{
		"thread_id": threadId,
		"user_id":   userId,
	}

	var messages []*Message
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var message Message
		err = cursor.Decode(&message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}
