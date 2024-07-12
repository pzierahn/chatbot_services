package datastore

import (
	"context"
	"github.com/google/uuid"
)

type Error struct {
	Id         uuid.UUID `bson:"_id,omitempty"`
	UserId     string    `bson:"user_id,omitempty"`
	Exception  string    `bson:"exception,omitempty"`
	StackTrace string    `bson:"stack_trace,omitempty"`
	AppVersion string    `bson:"app_version,omitempty"`
}

func (service *Service) InsertError(ctx context.Context, error *Error) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionCrashlytics)

	_, err := coll.InsertOne(ctx, error)
	if err != nil {
		return err
	}

	return nil
}
