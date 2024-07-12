package datastore

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Payments struct {
	Id     uuid.UUID `bson:"_id,omitempty"`
	UserId string    `bson:"user_id,omitempty"`
	Amount int       `bson:"amount,omitempty"`
	Date   time.Time `bson:"date,omitempty"`
}

func (service *Service) InsertPayment(ctx context.Context, payment *Payments) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionPayments)

	_, err := coll.InsertOne(ctx, payment)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) GetPayments(ctx context.Context, userId string) ([]Payments, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionPayments)

	cur, err := coll.Find(ctx, map[string]interface{}{
		"user_id": userId,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()

	var payments []Payments
	err = cur.All(ctx, &payments)
	if err != nil {
		return nil, err
	}

	return payments, nil
}
