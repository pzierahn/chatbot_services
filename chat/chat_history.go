package chat

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type chatMessage struct {
	id           string
	userId       string
	collectionId string
	prompt       string
	completion   string
	references   []string
}

func (service *Service) storeThread(ctx context.Context, message chatMessage) error {

	transaction, err := service.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = transaction.Rollback(ctx) }()

	_, err = transaction.Exec(ctx,
		`INSERT INTO threads (id, user_id, collection_id, prompt, completion)
		VALUES ($1, $2, $3, $4, $5)`,
		message.id,
		message.userId,
		message.collectionId,
		message.prompt,
		message.completion)
	if err != nil {
		return err
	}

	for _, source := range message.references {
		_, err = transaction.Exec(ctx,
			`INSERT INTO thread_references (thread_id, document_chunk_id)
			VALUES ($1, $2)`,
			message.id, source)
		if err != nil {
			return err
		}
	}

	return transaction.Commit(ctx)
}
