package datastore

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

// Service defines the datastore service
type Service struct {
	mongo *mongo.Client
}

// Close closes the connection to the database
func (service *Service) Close() {
	_ = service.mongo.Disconnect(context.Background())
}

const (
	DatabaseName = "chatbot"
)

const (
	CollectionThreads      = "threads"
	CollectionCollections  = "collections"
	CollectionDokuments    = "documents"
	CollectionCrashlytics  = "crashlytics"
	CollectionPayments     = "payments"
	CollectionModelUsages  = "model_usages"
	CollectionNotionAPIKey = "notion_api_keys"
)

func New(ctx context.Context) (*Service, error) {
	uri := os.Getenv("CHATBOT_MONGODB_URI")
	if uri == "" {
		return nil, errors.New("CHATBOT_MONGODB_URI not set")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &Service{
		mongo: client,
	}, nil
}
