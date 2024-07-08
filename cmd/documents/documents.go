package main

import (
	"cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"github.com/pzierahn/chatbot_services/auth"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/documents"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/llm/anthropic"
	"github.com/pzierahn/chatbot_services/llm/openai"
	"github.com/pzierahn/chatbot_services/llm/vertex"
	"github.com/pzierahn/chatbot_services/vectordb"
	"github.com/pzierahn/chatbot_services/vectordb/qdrant"
	"google.golang.org/api/option"
	"log"
	"os"
)

const credentialsFile = "service_account.json"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func initDatastore(ctx context.Context) *datastore.Service {
	uri := os.Getenv("CHATBOT_MONGODB_URI")
	if uri == "" {
		log.Fatal("CHATBOT_MONGODB_URI is not set")
	}

	db, err := datastore.New(ctx)
	if err != nil {
		log.Fatalf("failed to create datastore service: %v", err)
	}

	return db
}

func initBucket(ctx context.Context) *storage.BucketHandle {
	var opts []option.ClientOption
	if _, err := os.Stat(credentialsFile); err == nil {
		serviceAccount := option.WithCredentialsFile(credentialsFile)
		opts = append(opts, serviceAccount)
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		log.Fatalf("failed to create firebase app: %v", err)
	}

	firebaseStorage, err := app.Storage(ctx)
	if err != nil {
		log.Fatalf("failed to create firebase storage client: %v", err)
	}

	bucket, err := firebaseStorage.Bucket("brainboost-399710.appspot.com")
	if err != nil {
		log.Fatalf("did not get bucket: %v", err)
	}

	return bucket
}

func initModels(ctx context.Context) []llm.Chat {
	openaiClient, err := openai.New()
	if err != nil {
		log.Fatalf("failed to create openai client: %v", err)
	}

	vertexClient, err := vertex.New(ctx)
	if err != nil {
		log.Fatalf("failed to create vertex client: %v", err)
	}

	claude, err := anthropic.New()
	if err != nil {
		log.Fatalf("failed to create anthropic client: %v", err)
	}

	models := []llm.Chat{
		openaiClient,
		vertexClient,
		claude,
	}

	return models
}

func initSearch() vectordb.DB {
	search, err := qdrant.New()
	if err != nil {
		log.Fatalf("failed to create qdrant search: %v", err)
	}

	return search
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	database := initDatastore(ctx)
	search := initSearch()
	bucket := initBucket(ctx)
	fakeAuth, _ := auth.WithInsecure()

	documentsService := &documents.Service{
		Auth:        fakeAuth,
		Database:    database,
		Storage:     bucket,
		SearchIndex: search,
	}
}
