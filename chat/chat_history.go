package chat

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type chatMessage struct {
	userID       uuid.UUID
	collectionID string
	prompt       string
	completion   string
	references   []uuid.UUID
}

func (service *Service) storeChatMessage(ctx context.Context, message chatMessage) (uuid.UUID, error) {

	transaction, err := service.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return uuid.Nil, err
	}
	defer func() { _ = transaction.Rollback(ctx) }()

	var chatID uuid.UUID
	err = transaction.QueryRow(ctx,
		`INSERT INTO chat_message (user_id, collection_id, prompt, completion)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		message.userID,
		message.collectionID,
		message.prompt,
		message.completion).Scan(&chatID)
	if err != nil {
		return uuid.Nil, err
	}

	for _, source := range message.references {
		_, err = transaction.Exec(ctx,
			`INSERT INTO chat_message_source (chat_message_id, document_embeddings_id)
			VALUES ($1, $2)`,
			chatID, source)
		if err != nil {
			return uuid.Nil, err
		}
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	return chatID, nil
}

func (service *Service) GetChatMessages(ctx context.Context, collection *pb.Collection) (*pb.ChatMessages, error) {
	uid, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := service.db.Query(ctx,
		`SELECT id FROM chat_message
          WHERE user_id = $1 AND collection_id = $2
          ORDER BY created_at DESC`,
		uid, collection)
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

func (service *Service) getChatMessageDocuments(ctx context.Context, userID uuid.UUID, message *pb.ChatMessage) ([]*pb.ChatMessage_Document, error) {
	rows, err := service.db.Query(ctx,
		`SELECT de.id, doc.id, doc.filename, de.page
			FROM chat_message AS cm,
				 chat_message_source AS cms,
				 documents AS doc,
				 document_embeddings as de
			WHERE cm.id = cms.chat_message_id AND
			      cms.document_embeddings_id = de.id AND
			      doc.id = de.document_id AND
			      cm.id = $1 AND
			      cm.user_id = $2 AND
			      doc.collection_id = $3
			ORDER BY de.page`,
		message.Id, userID, message.CollectionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		docIDs = make(map[uuid.UUID]string)
		pages  = make(map[uuid.UUID][]uint32)
	)

	for rows.Next() {
		var (
			id       uuid.UUID
			docID    uuid.UUID
			filename string
			page     uint32
		)

		if err = rows.Scan(
			&id,
			&docID,
			&filename,
			&page); err != nil {
			return nil, err
		}

		docIDs[id] = filename
		pages[id] = append(pages[id], page)
	}

	var docs []*pb.ChatMessage_Document
	for docID := range docIDs {
		docs = append(docs, &pb.ChatMessage_Document{
			Id:       docID.String(),
			Filename: docIDs[docID],
			Pages:    pages[docID],
		})
	}

	return docs, nil
}

func (service *Service) GetChatMessage(ctx context.Context, id *pb.MessageID) (*pb.ChatMessage, error) {
	userID, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	var message pb.ChatMessage
	var createdAt time.Time

	err = service.db.QueryRow(ctx,
		`SELECT id, collection_id, created_at, prompt, completion
			FROM chat_message
			WHERE id = $1 AND
			      user_id = $2`,
		id, userID).Scan(
		&message.Id,
		&message.CollectionId,
		&createdAt,
		&message.Prompt,
		&message.Text)
	if err != nil {
		return nil, err
	}

	message.Timestamp = timestamppb.New(createdAt)

	message.Documents, err = service.getChatMessageDocuments(ctx, userID, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
