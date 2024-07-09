package migration

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"log"
	"time"
)

func (migrator *Migrator) MigratePayments() {
	ctx := context.Background()

	log.Printf("Migrating collections")

	rows, err := migrator.Legacy.Query(ctx, "SELECT id, user_id, date, amount FROM payments")
	if err != nil {
		log.Fatalf("Query collections: %v", err)
	}

	count := 0
	for rows.Next() {
		var (
			id     uuid.UUID
			userId string
			date   time.Time
			amount int
		)

		err = rows.Scan(&id, &userId, &date, &amount)
		if err != nil {
			log.Fatal(err)
		}

		// Insert collection into the new database
		err = migrator.Next.InsertPayment(ctx, &datastore.Payments{
			Id:     id,
			UserId: userId,
			Date:   date,
			Amount: amount,
		})
		if err != nil {
			log.Fatal(err)
		}

		count++
	}

	log.Printf("Migrated %d payments", count)
}
