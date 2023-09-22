package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

type ChatMessage struct {
	ID         uuid.UUID
	UID        uuid.UUID
	CreateAt   *time.Time
	Collection uuid.UUID
	Prompt     string
	Completion string
	Sources    []ChatMessageSource
}

type ChatMessageSource struct {
	ID           *uuid.UUID
	DocumentPage uuid.UUID
	Filename     string
	Page         int
}

func (client *Client) CreateChat(ctx context.Context, history ChatMessage) (*uuid.UUID, error) {
	transaction, err := client.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = transaction.Rollback(ctx) }()

	id := uuid.New()

	_, err = transaction.Exec(ctx,
		`INSERT INTO chat_message (id, uid, collection_id, prompt, completion)
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
			`INSERT INTO chat_message_source (chat_message_id, document_embeddings_id)
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

func (client *Client) GetChatMessages(ctx context.Context, uid, collection string) (ids []uuid.UUID, _ error) {
	rows, err := client.conn.Query(ctx,
		`SELECT id FROM chat_message
          WHERE uid = $1 AND collection_id = $2
          ORDER BY created_at DESC`,
		uid, collection)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (client *Client) GetChatMessage(ctx context.Context, id uuid.UUID, uid string) (*ChatMessage, error) {
	var message ChatMessage

	err := client.conn.QueryRow(ctx,
		`SELECT id, created_at, prompt, completion
			FROM chat_message
			WHERE id = $1 AND uid = $2`,
		id, uid).Scan(
		&message.ID,
		&message.CreateAt,
		&message.Prompt,
		&message.Completion)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	message.Sources, err = client.GetChatMessageDocuments(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (client *Client) GetChatMessageDocuments(ctx context.Context, id uuid.UUID, uid string) ([]ChatMessageSource, error) {
	rows, err := client.conn.Query(ctx,
		`SELECT de.id, doc.id, doc.filename, de.page
			FROM chat_message AS cm,
				 chat_message_source AS cms,
				 documents AS doc,
				 document_embeddings as de
			WHERE cm.id = cms.chat_message_id
			  AND cms.document_embeddings_id = de.id
			  AND doc.id = de.document_id
			  AND cm.id = $1
			  AND cm.uid = $2`,
		id, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []ChatMessageSource
	for rows.Next() {
		var source ChatMessageSource
		if err = rows.Scan(
			&source.ID,
			&source.DocumentPage,
			&source.Filename,
			&source.Page); err != nil {
			return nil, err
		}

		sources = append(sources, source)
	}

	return sources, nil
}
