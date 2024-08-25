package main

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm/openai"
	"github.com/pzierahn/chatbot_services/migration"
	pinecone_search "github.com/pzierahn/chatbot_services/search/pinecone"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	uri := os.Getenv("CHATBOT_MONGODB_URI")
	next, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}

	engine, err := openai.New()
	if err != nil {
		log.Fatalf("failed to create openai service: %v", err)
	}

	index, err := pinecone_search.New(engine, "documents")
	if err != nil {
		log.Fatalf("failed to create search service: %v", err)
	}

	// Create a new migrator
	migrator := &migration.Migrator{
		Search:   index,
		Database: next,
	}

	migrator.MigrateVectorDB()
}
