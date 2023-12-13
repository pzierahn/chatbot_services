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

		// Insert into new database
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

func MigrateOpenaiUsage(from, to *pgxpool.Pool) {

	// Get user mapping
	userMapping := GetUserIdMapping()

	ctx := context.Background()

	// Get all payments from supabase

	rows, err := from.Query(ctx, `
		SELECT id, user_id, created_at, model, input_tokens, output_tokens
		FROM openai_usage`)
	if err != nil {
		log.Fatalln(err)
	}

	// Iterate over payments
	for rows.Next() {
		var (
			id     string
			userId string
			date   time.Time
			model  string
			input  string
			output string
		)

		err = rows.Scan(&id, &userId, &date, &model, &input, &output)
		if err != nil {
			log.Fatalln(err)
		}

		newUserId, ok := userMapping[userId]
		if !ok {
			log.Fatalf("user not found: %v", userId)
		}

		log.Printf("Usage: %v", id)

		_, err = to.Exec(ctx, `
			INSERT INTO openai_usages (id, user_id, created_at, model, input_tokens, output_tokens)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT DO NOTHING
		`, id, newUserId, date, model, input, output)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
