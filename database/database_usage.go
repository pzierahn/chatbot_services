package database

import (
	"context"
	"github.com/google/uuid"
	"time"
)

// Usage represents a row in the openai_usage table
type Usage struct {
	ID        uuid.UUID
	CreatedAt *time.Time
	UID       string
	Model     string
	Input     int
	Output    int
}

// CreateUsage inserts a new usage record into the openai_usage table
func (client *Client) CreateUsage(ctx context.Context, usage Usage) error {
	_, err := client.conn.Exec(ctx,
		`INSERT INTO openai_usage (id, uid, model, input, output)
			VALUES ($1, $2, $3, $4, $5)`,
		usage.ID, usage.UID, usage.Model, usage.Input, usage.Output)
	return err
}

// ReadUsage retrieves a usage record by ID from the openai_usage table
func (client *Client) ReadUsage(ctx context.Context, id string) (Usage, error) {
	var usage Usage
	err := client.conn.QueryRow(ctx,
		`SELECT id, created_at, uid, model, input, output
			FROM openai_usage WHERE id = $1`, id).
		Scan(&usage.ID, &usage.CreatedAt, &usage.UID, &usage.Model, &usage.Input, &usage.Output)
	return usage, err
}

// UpdateUsage updates an existing usage record in the openai_usage table
func (client *Client) UpdateUsage(ctx context.Context, usage Usage) error {
	_, err := client.conn.Exec(ctx,
		`UPDATE openai_usage
			SET uid = $2, model = $3, input = $4, output = $5
			WHERE id = $1`,
		usage.ID, usage.UID, usage.Model, usage.Input, usage.Output)
	return err
}

// DeleteUsage deletes a usage record by ID from the openai_usage table
func (client *Client) DeleteUsage(ctx context.Context, id string) error {
	_, err := client.conn.Exec(ctx, `DELETE FROM openai_usage WHERE id = $1`, id)
	return err
}
