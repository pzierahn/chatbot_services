package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type ChatMessage struct {
	ID         *uuid.UUID
	UID        string
	CreateAt   *time.Time
	Collection string
	Prompt     string
	Completion string
	Sources    []ChatMessageSource
}

type ChatMessageSource struct {
	ID           *uuid.UUID
	DocumentPage uuid.UUID
}

func (client *Client) CreateChat(ctx context.Context, history ChatMessage) (*uuid.UUID, error) {
	transaction, err := client.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = transaction.Rollback(ctx) }()

	id := uuid.New()

	_, err = transaction.Exec(ctx,
		`INSERT INTO chat_message (id, uid, collection, prompt, completion)
		VALUES ($1, $2, $3, $4, $5)`,
		id,
		history.UID,
		history.Collection,
		history.Prompt,
		history.Completion)
	if err != nil {
		return nil, err
	}

	for _, source := range history.Sources {
		err = transaction.QueryRow(ctx,
			`INSERT INTO chat_message_source (chat, document_page)
			VALUES ($1, $2)
			RETURNING id`,
			id,
			source.DocumentPage).
			Scan(&source.ID)
		if err != nil {
			return nil, err
		}
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
