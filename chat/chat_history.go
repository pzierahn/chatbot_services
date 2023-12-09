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
	id           string
	userId       string
	collectionId string
	prompt       string
	completion   string
	references   []uuid.UUID
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
	uid, err := service.auth.ValidateToken(ctx)
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

func (service *Service) getChatMessageDocuments(ctx context.Context, userId string, message *pb.ChatMessage) ([]*pb.ChatMessage_Document, error) {
	rows, err := service.db.Query(ctx,
		`SELECT de.id, doc.id, doc.filename, de.page
			FROM chat_messages AS cm,
				 chat_message_references AS cms,
				 documents AS doc,
				 document_chunks as de
			WHERE cm.id = cms.chat_message_id AND
			      cms.document_chunk_id = de.id AND
			      doc.id = de.document_id AND
			      cm.id = $1 AND
			      cm.user_id = $2 AND
			      doc.collection_id = $3
			ORDER BY de.page`,
		message.Id, userId, message.CollectionId)
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
	userId, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	var message pb.ChatMessage

	var pompt string
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
		&pompt,
		&message.Text)
	if err != nil {
		return nil, err
	}

	message.Prompt = &pb.Prompt{
		Prompt:       pompt,
		CollectionId: message.CollectionId,
	}
	message.Timestamp = timestamppb.New(createdAt)

	message.Documents, err = service.getChatMessageDocuments(ctx, userId, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
