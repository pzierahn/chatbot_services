package datastore

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Collection struct {
	// ID of the collection
	Id uuid.UUID `bson:"_id,omitempty"`

	// User ID
	UserId string `bson:"user_id,omitempty"`

	// Name of the collection
	Name string `bson:"name,omitempty"`
}

// StoreCollection stores a collection in the database
func (service *Service) StoreCollection(ctx context.Context, collection *Collection) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionCollections)

	_, err := coll.InsertOne(ctx, collection)
	if err != nil {
		return err
	}

	return nil
}

// GetCollection retrieves a collection from the database
func (service *Service) GetCollection(ctx context.Context, userId string, collectionId uuid.UUID) (*Collection, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionCollections)

	var collection Collection
	err := coll.FindOne(ctx, bson.M{
		"_id":     collectionId,
		"user_id": userId,
	}).Decode(&collection)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

// GetCollections retrieves all collections from the database
func (service *Service) GetCollections(ctx context.Context, userId string) ([]Collection, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionCollections)

	cursor, err := coll.Find(ctx, bson.M{
		"user_id": userId,
	})
	if err != nil {
		return nil, err
	}

	var collections []Collection
	err = cursor.All(ctx, &collections)
	if err != nil {
		return nil, err
	}

	return collections, nil
}
