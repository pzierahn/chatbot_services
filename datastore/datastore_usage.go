package datastore

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type ModelUsage struct {
	Id           uuid.UUID `bson:"_id,omitempty"`
	UserId       string    `bson:"user_id,omitempty"`
	Timestamp    time.Time `bson:"timestamp,omitempty"`
	ModelId      string    `bson:"model_id,omitempty"`
	InputTokens  uint32    `bson:"input_tokens,omitempty"`
	OutputTokens uint32    `bson:"output_tokens,omitempty"`
}

func (service *Service) InsertModelUsage(ctx context.Context, usage *ModelUsage) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionModelUsages)

	_, err := coll.InsertOne(ctx, usage)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) GetModelUsages(ctx context.Context, userId string) ([]ModelUsage, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionModelUsages)

	cur, err := coll.Find(ctx, map[string]interface{}{
		"user_id": userId,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()

	var usages []ModelUsage
	err = cur.All(ctx, &usages)
	if err != nil {
		return nil, err
	}

	return usages, nil
}
