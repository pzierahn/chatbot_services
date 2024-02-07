package chat

import (
	"context"
	"github.com/jackc/pgx/v5"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/utils"
)

func (service *Service) storeThread(ctx context.Context, userId, collectionId string, thread *pb.Thread) error {

	tx, err := service.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx,
		`INSERT INTO threads (id, user_id, collection_id)
			VALUES ($1, $2, $3)`,
		thread.Id,
		userId,
		collectionId,
	)
	if err != nil {
		return err
	}

	for _, msg := range thread.Messages {
		_, err = tx.Exec(ctx,
			`INSERT INTO thread_messages (id, user_id, thread_id, created_at, prompt, completion)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			msg.Id,
			userId,
			thread.Id,
			utils.ProtoToTime(msg.Timestamp),
			msg.Prompt,
			msg.Completion)
		if err != nil {
			return err
		}
	}

	for _, source := range thread.ReferenceIDs {
		_, err = tx.Exec(ctx,
			`INSERT INTO thread_references (user_id, thread_id, document_chunk_id)
			VALUES ($1, $2, $3)`,
			userId,
			thread.Id,
			source)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
