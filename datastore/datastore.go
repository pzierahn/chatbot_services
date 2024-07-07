package datastore

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	DatabaseName       = "chatbot"
	CollectionMessages = "messages"
)

func New(ctx context.Context) (*Service, error) {
	//uri := os.Getenv("CHATBOT_MONGODB_URI")
	//if uri == "" {
	//	log.Fatal("MONGODB_URI is not set")
	//}

	uri := "mongodb://localhost:27017"

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &Service{
		mongo: client,
	}, nil
}
