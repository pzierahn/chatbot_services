package migration

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

// UpdateCollections updates the user_id column in the collections table.
func UpdateCollections(ctx context.Context, db *pgxpool.Pool) {

	// Get user mapping
	userMapping := GetUserIdMapping()

	// Update collections
	for oldId, newId := range userMapping {
		log.Printf("Update collection: %v -> %v", oldId, newId)

		_, err := db.Exec(ctx,
			`UPDATE collections
				SET user_id = $1
				WHERE id = $2`,
			newId, oldId)
		if err != nil {
			log.Fatalf("did not update: %v", err)
		}
	}
}
