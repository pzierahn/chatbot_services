package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/llm/openai"
	"github.com/pzierahn/chatbot_services/migration"
	"github.com/pzierahn/chatbot_services/search/qdrant"
	"log"
	"os"
)

func init() {
	_ = os.Setenv("CHATBOT_MONGODB_URI", "mongodb://localhost:27017")
	_ = os.Setenv("CHATBOT_QDRANT_KEY", "")
	_ = os.Setenv("CHATBOT_QDRANT_URL", "localhost:6334")
	_ = os.Setenv("CHATBOT_QDRANT_INSECURE", "true")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	addr := os.Getenv("CHATBOT_DB")
	legacy, err := pgxpool.New(ctx, addr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer legacy.Close()

	next, err := datastore.New(ctx)
	if err != nil {
		log.Fatalf("failed to create datastore service: %v", err)
	}
	defer next.Close()

	engine, err := openai.New()
	if err != nil {
		log.Fatalf("failed to create openai service: %v", err)
	}

	index, err := qdrant.New(engine, "documents_v2")
	if err != nil {
		log.Fatalf("failed to create search service: %v", err)
	}

	// Create a new migrator
	migrator := &migration.Migrator{
		Legacy: legacy,
		Next:   next,
	}

	//migrator.MigrateCollections()
	//migrator.MigrateDocuments()
	migrator.MigrateDocumentToSearch(index)
	//migrator.MigrateThreads()
	//migrator.MigratePayments()
	//migrator.MigrateUsages()
}
