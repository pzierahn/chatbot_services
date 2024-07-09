package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/migration"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	addr := os.Getenv("CHATBOT_DB")
	legacy, err := pgxpool.New(ctx, addr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer legacy.Close()

	uri := "mongodb://localhost:27017"
	next, err := datastore.NewFrom(ctx, uri)
	if err != nil {
		log.Fatalf("failed to create datastore service: %v", err)
	}
	defer next.Close()

	// Create a new migrator
	migrator := &migration.Migrator{
		Legacy: legacy,
		Next:   next,
	}

	//migrator.MigrateCollections()
	//migrator.MigrateDocuments()
	migrator.MigrateThreads()
}
