package database

import (
	"context"
	"github.com/google/uuid"
	"time"
)

// Usage represents a row in the openai_usage table
type Usage struct {
	ID        uuid.UUID
	UID       uuid.UUID
	CreatedAt *time.Time
	Model     string
	Input     uint32
	Output    uint32
}

type ModelUsage struct {
	Model  string
	Input  uint32
	Output uint32
}

// CreateUsage inserts a new usage record into the openai_usage table
func (client *Client) CreateUsage(ctx context.Context, usage Usage) (uuid.UUID, error) {
	err := client.conn.QueryRow(ctx,
		`INSERT INTO openai_usage (user_id, model, input, output)
			VALUES ($1, $2, $3, $4)
			RETURNING id`,
		usage.UID, usage.Model, usage.Input, usage.Output).
		Scan(&usage.ID)

	return usage.ID, err
}

// GetModelUsages retrieves a usage record by ID from the openai_usage table
func (client *Client) GetModelUsages(ctx context.Context, uid uuid.UUID) ([]ModelUsage, error) {
	rows, err := client.conn.Query(ctx,
		`SELECT model, SUM(input), SUM(output)
			FROM openai_usage
			WHERE user_id = $1
			GROUP BY model`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usages []ModelUsage
	for rows.Next() {
		var usage ModelUsage
		err = rows.Scan(&usage.Model, &usage.Input, &usage.Output)
		if err != nil {
			return nil, err
		}
		usages = append(usages, usage)
	}

	return usages, nil
}
