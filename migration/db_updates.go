package migration

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

func MigratePayments(from, to *pgxpool.Pool) {

	// Get user mapping
	userMapping := GetUserIdMapping()

	ctx := context.Background()

	// Get all payments from supabase

	rows, err := from.Query(ctx, `SELECT id, user_id, date, amount FROM payments`)
	if err != nil {
		log.Fatalln(err)
	}

	// Iterate over payments
	for rows.Next() {
		var (
			id     string
			userId string
			date   time.Time
			amount string
		)

		err = rows.Scan(&id, &userId, &date, &amount)
		if err != nil {
			log.Fatalln(err)
		}

		newUserId, ok := userMapping[userId]
		if !ok {
			log.Fatalf("user not found: %v", userId)
		}

		log.Printf("Payment: %v", id)

		_, err = to.Exec(ctx, `
			INSERT INTO payments (id, user_id, date, amount)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT DO NOTHING
		`, id, newUserId, date, amount)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
