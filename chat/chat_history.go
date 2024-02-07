package chat

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type chatMessage struct {
	id           string
	userId       string
	collectionId string
	prompt       string
	completion   string
	references   []string
}

func (service *Service) storeChatMessage(ctx context.Context, message chatMessage) error {

	transaction, err := service.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = transaction.Rollback(ctx) }()

	_, err = transaction.Exec(ctx,
		`INSERT INTO chat_messages (id, user_id, collection_id, prompt, completion)
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
			`INSERT INTO chat_message_references (chat_message_id, document_chunk_id)
			VALUES ($1, $2)`,
			message.id, source)
		if err != nil {
			return err
		}
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) GetChatMessages(ctx context.Context, collection *pb.Collection) (*pb.ChatMessages, error) {
	uid, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := service.db.Query(ctx,
		`SELECT id FROM chat_messages
          WHERE user_id = $1 AND
                collection_id = $2
          ORDER BY created_at DESC`,
		uid, collection.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := &pb.ChatMessages{}
	for rows.Next() {
		var id uuid.UUID
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		messages.Ids = append(messages.Ids, id.String())
	}

	return messages, nil
}

func (service *Service) getChatMessageReferences(ctx context.Context, messageId string) ([]string, error) {
	rows, err := service.db.Query(ctx,
		`SELECT document_chunk_id
			FROM chat_message_references
			WHERE chat_message_id = $1`,
		messageId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var references []string
	for rows.Next() {
		var id string

		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		references = append(references, id)
	}

	return references, nil
}

func (service *Service) GetChatMessage(ctx context.Context, id *pb.MessageID) (*pb.ChatMessage, error) {
	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	var message pb.ChatMessage

	var createdAt time.Time

	err = service.db.QueryRow(ctx,
		`SELECT id, collection_id, created_at, prompt, completion
			FROM chat_messages
			WHERE id = $1 AND
			      user_id = $2`,
		id.Id, userId).Scan(
		&message.Id,
		&message.CollectionId,
		&createdAt,
		&message.Prompt,
		&message.Text)
	if err != nil {
		return nil, err
	}

	message.Timestamp = timestamppb.New(createdAt)

	message.References, err = service.getChatMessageReferences(ctx, message.Id)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
