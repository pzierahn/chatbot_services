package datastore

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type NotionAPIKey struct {
	Id     uuid.UUID `bson:"_id,omitempty"`
	UserId string    `bson:"user_id,omitempty"`
	Key    string    `bson:"api_key,omitempty"`
}

func (service *Service) InsertNotionAPIKey(ctx context.Context, key *NotionAPIKey) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionNotionAPIKey)
	_, err := coll.InsertOne(ctx, key)
	return err
}

func (service *Service) UpdateNotionAPIKey(ctx context.Context, key *NotionAPIKey) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionNotionAPIKey)

	_, err := coll.UpdateOne(ctx, bson.M{
		"user_id": key.UserId,
	}, bson.M{
		"$set": bson.M{
			"api_key": key.Key,
		},
	})
	return err
}

func (service *Service) DeleteNotionAPIKey(ctx context.Context, userId string) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionNotionAPIKey)
	_, err := coll.DeleteOne(ctx, bson.M{
		"user_id": userId,
	})
	return err
}

func (service *Service) GetNotionAPIKey(ctx context.Context, userId string) (string, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionNotionAPIKey)

	var key NotionAPIKey
	err := coll.FindOne(ctx, bson.M{
		"user_id": userId,
	}).Decode(&key)
	if err != nil {
		return "", err
	}

	return key.Key, nil
}
