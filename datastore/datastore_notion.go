package datastore

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type NotionAPIKey struct {
	UserId string `bson:"_id,omitempty"`
	Key    string `bson:"api_key,omitempty"`
}

func (service *Service) InsertNotionAPIKey(ctx context.Context, key *NotionAPIKey) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionNotionAPIKey)
	_, err := coll.InsertOne(ctx, key)
	return err
}

func (service *Service) UpdateNotionAPIKey(ctx context.Context, key *NotionAPIKey) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionNotionAPIKey)

	_, err := coll.UpdateOne(ctx, bson.M{
		"_id": key.UserId,
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
		"_id": userId,
	})
	return err
}

func (service *Service) GetNotionAPIKey(ctx context.Context, userId string) (string, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionNotionAPIKey)

	var key NotionAPIKey
	err := coll.FindOne(ctx, bson.M{
		"_id": userId,
	}).Decode(&key)
	if err != nil {
		return "", err
	}

	return key.Key, nil
}
