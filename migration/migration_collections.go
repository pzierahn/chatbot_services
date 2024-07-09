package migration

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"log"
)

func (migrator *Migrator) MigrateCollections() {
	ctx := context.Background()

	log.Printf("Migrating collections")

	// Query all collections from the legacy database
	rows, err := migrator.Legacy.Query(ctx, "SELECT id, user_id, name FROM collections")
	if err != nil {
		log.Fatalf("Query collections: %v", err)
	}

	// Iterate over all collections
	for rows.Next() {
		var id uuid.UUID
		var userId, name string
		err = rows.Scan(&id, &userId, &name)
		if err != nil {
			log.Fatalf("Scan collection: %v", err)
		}

		// Insert collection into the new database
		err = migrator.Next.InsertCollection(ctx, &datastore.Collection{
			Id:     id,
			UserId: userId,
			Name:   name,
		})
		if err != nil {
			log.Fatalf("Insert collection: %v", err)
		}
	}
}
