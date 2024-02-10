package chat

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type reference struct {
	Id         string
	DocumentId string
	Page       int32
	Text       string
}

func (service *Service) getReferences(ctx context.Context, uid, threadId string) ([]*reference, error) {
	rows, err := service.db.Query(
		ctx,
		`SELECT dc.id, dc.document_id, dc.page, text
			FROM thread_references as tr, document_chunks as dc
			WHERE user_id = $1 AND 
			      thread_id = $2 AND
			      tr.document_chunk_id = dc.id
		  	ORDER BY dc.document_id, dc.page`,
		uid, threadId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []*reference
	for rows.Next() {
		var ref reference

		err = rows.Scan(&ref.Id, &ref.DocumentId, &ref.Page, &ref.Text)
		if err != nil {
			return nil, err
		}

		refs = append(refs, &ref)
	}

	return refs, nil
}

func (service *Service) getThreadMessages(ctx context.Context, uid, threadId string) ([]*pb.Message, error) {
	rows, err := service.db.Query(
		ctx,
		`SELECT id, prompt, completion, created_at
			FROM thread_messages
			WHERE user_id = $1 AND thread_id = $2
			ORDER BY created_at`,
		uid, threadId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*pb.Message

	for rows.Next() {
		var message pb.Message
		var createdAt time.Time

		err = rows.Scan(&message.Id, &message.Prompt, &message.Completion, &createdAt)
		if err != nil {
			return nil, err
		}

		message.Timestamp = timestamppb.New(createdAt)
		messages = append(messages, &message)
	}

	return messages, nil
}

func (service *Service) GetThread(ctx context.Context, thread *pb.ThreadID) (*pb.Thread, error) {
	userId, err := service.Verify(ctx)
	if err != nil {
		return nil, err
	}

	references, err := service.getReferences(ctx, userId, thread.Id)
	if err != nil {
		return nil, err
	}

	referenceIds := make([]string, len(references))
	for inx, ref := range references {
		referenceIds[inx] = ref.Id
	}

	massages, err := service.getThreadMessages(ctx, userId, thread.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Thread{
		Id:           thread.Id,
		ReferenceIDs: referenceIds,
		Messages:     massages,
	}, nil
}
