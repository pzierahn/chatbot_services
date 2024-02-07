package chat

import (
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"strings"
)

func (service *Service) getThreadMessages(ctx context.Context, uid, threadId string) ([]*llm.Message, error) {
	rows, err := service.db.Query(
		ctx,
		`SELECT prompt, completion
			FROM thread_messages
			WHERE user_id = $1 AND thread_id = $2
			ORDER BY created_at`,
		uid, threadId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*llm.Message

	for rows.Next() {
		var prompt, completion string
		err = rows.Scan(&prompt, &completion)
		if err != nil {
			return nil, err
		}

		messages = append(messages, []*llm.Message{{
			Type: llm.MessageTypeUser,
			Text: prompt,
		}, {
			Type: llm.MessageTypeBot,
			Text: completion,
		}}...)
	}

	return messages, nil
}

func (service *Service) getReferences(ctx context.Context, uid, threadId string) ([]*llm.Message, error) {
	rows, err := service.db.Query(
		ctx,
		`SELECT dc.id, text
			FROM thread_references as tr, document_chunks as dc
			WHERE user_id = $1 AND 
			      thread_id = $2 AND
			      tr.document_chunk_id = dc.id
		  	ORDER BY dc.page`,
		uid, threadId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// DocumentId --> pages
	docRefs := make(map[string][]string)
	for rows.Next() {
		var docId, text string
		err = rows.Scan(&text)
		if err != nil {
			return nil, err
		}

		docRefs[docId] = append(docRefs[docId], text)
	}

	var messages []*llm.Message
	for _, pages := range docRefs {
		messages = append(messages, &llm.Message{
			Type: llm.MessageTypeUser,
			Text: strings.Join(pages, "\n"),
		})
	}

	return messages, nil
}

func (service *Service) PostMessage(ctx context.Context, prompt *pb.Prompt) (*pb.Message, error) {
	userId, err := service.Verify(ctx)
	if err != nil {
		return nil, err
	}

	if prompt.ModelOptions == nil {
		return nil, fmt.Errorf("options missing")
	}

	var messages []*llm.Message

	references, err := service.getReferences(ctx, userId, prompt.ThreadID)
	if err != nil {
		log.Printf("Error fetching references: %v", err)
		return nil, err
	}
	messages = append(messages, references...)

	chat, err := service.getThreadMessages(ctx, userId, prompt.ThreadID)
	if err != nil {
		log.Printf("Error fetching thread messages: %v", err)
		return nil, err
	}
	messages = append(messages, chat...)

	messages = append(messages, &llm.Message{
		Type: llm.MessageTypeUser,
		Text: prompt.Prompt,
	})

	model, err := service.getModel(prompt.ModelOptions.Model)
	if err != nil {
		log.Printf("Error fetching model: %v", err)
		return nil, err
	}

	completion, err := model.GenerateCompletion(ctx, &llm.GenerateRequest{
		Messages:    messages,
		Model:       prompt.ModelOptions.Model,
		MaxTokens:   1024,
		Temperature: prompt.ModelOptions.Temperature,
		UserId:      userId,
	})
	if err != nil {
		log.Printf("Error generating completion: %v", err)
		return nil, err
	}

	message := &pb.Message{
		Prompt:     prompt.Prompt,
		Completion: completion.Text,
		Timestamp:  timestamppb.Now(),
	}

	err = service.db.QueryRow(
		ctx,
		`INSERT INTO thread_messages (user_id, thread_id, prompt, completion, created_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING (id)`,
		userId,
		prompt.ThreadID,
		prompt.Prompt,
		completion.Text,
		utils.ProtoToTime(message.Timestamp)).
		Scan(&message.Id)
	if err != nil {
		log.Printf("Error inserting message: %v", err)
		return nil, err
	}

	return message, err
}
