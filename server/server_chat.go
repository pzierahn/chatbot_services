package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/database"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
	"log"
	"sort"
	"strings"
)

type chatContext struct {
	fragments []string
	docs      []*pb.ChatMessage_Document
	pageIDs   []uuid.UUID
}

func (server *Server) getBackgroundFromPrompt(ctx context.Context, uid uuid.UUID, prompt *pb.Prompt) (*chatContext, error) {
	sort.Slice(prompt.Documents, func(i, j int) bool {
		return prompt.Documents[i].Filename < prompt.Documents[j].Filename
	})

	for _, doc := range prompt.Documents {
		sort.Slice(doc.Pages, func(i, j int) bool {
			return doc.Pages[i] < doc.Pages[j]
		})
	}

	bg := chatContext{}

	for _, doc := range prompt.Documents {
		id, err := uuid.Parse(doc.Id)
		if err != nil {
			return nil, err
		}

		content, err := server.db.GetPageContent(ctx, database.PageContentQuery{
			Id:     id,
			UserId: uid.String(),
			Pages:  doc.Pages,
		})
		if err != nil {
			return nil, err
		}

		var parts []string
		for _, page := range content {
			parts = append(parts, page.Text)
			bg.pageIDs = append(bg.pageIDs, page.Id)
		}

		bg.fragments = append(bg.fragments, strings.Join(parts, "\n"))
		bg.docs = append(bg.docs, &pb.ChatMessage_Document{
			Id:       doc.Id,
			Filename: doc.Filename,
			Pages:    doc.Pages,
		})
	}

	return &bg, nil
}

func (server *Server) getBackgroundFromDB(ctx context.Context, uid uuid.UUID, prompt *pb.Prompt) (*chatContext, error) {

	collection, err := uuid.Parse(prompt.Collection)
	if err != nil {
		return nil, err
	}

	query := SearchQuery{
		UserId:     uid,
		Collection: collection,
		Prompt:     prompt.Prompt,
		Limit:      int(prompt.Options.Limit),
		Threshold:  prompt.Options.Threshold,
	}

	results, err := server.SearchDocuments(ctx, query)
	if err != nil {
		return nil, err
	}

	bg := chatContext{}

	for _, doc := range results {
		text := make([]string, len(doc.Pages))
		pages := make([]uint32, len(doc.Pages))
		scores := make([]float32, len(doc.Pages))

		for iny, page := range doc.Pages {
			text[iny] = page.Text
			pages[iny] = page.Page
			scores[iny] = page.Score

			bg.pageIDs = append(bg.pageIDs, page.Id)
		}

		bg.fragments = append(bg.fragments, strings.Join(text, "\n"))
		bg.docs = append(bg.docs, &pb.ChatMessage_Document{
			Id:       doc.DocId.String(),
			Filename: doc.Filename,
			Pages:    pages,
			Scores:   scores,
		})
	}

	return &bg, nil
}

func (server *Server) Chat(ctx context.Context, prompt *pb.Prompt) (*pb.ChatMessage, error) {
	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	byt, _ := json.MarshalIndent(prompt, "", "  ")

	log.Printf("Chat: %s", byt)

	if prompt.Options == nil {
		return nil, fmt.Errorf("options missing")
	}

	var bg *chatContext
	if prompt.Documents == nil || len(prompt.Documents) == 0 {
		bg, err = server.getBackgroundFromDB(ctx, uid, prompt)
	} else {
		bg, err = server.getBackgroundFromPrompt(ctx, uid, prompt)
	}

	if err != nil {
		return nil, err
	}

	var messages []openai.ChatCompletionMessage
	for _, text := range bg.fragments {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt.Prompt,
	})

	resp, err := server.gpt.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       prompt.Options.Model,
			Temperature: prompt.Options.Temperature,
			MaxTokens:   int(prompt.Options.MaxTokens),
			Messages:    messages,
			N:           1,
			User:        uid.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	_, err = server.db.CreateUsage(ctx, database.Usage{
		UID:    uid,
		Model:  resp.Model,
		Input:  uint32(resp.Usage.PromptTokens),
		Output: uint32(resp.Usage.CompletionTokens),
	})
	if err != nil {
		log.Printf("Chat: error %v", err)
	}

	completion := &pb.ChatMessage{
		Prompt:    prompt,
		Text:      resp.Choices[0].Message.Content,
		Documents: bg.docs,
	}

	server.storeChatMessage(ctx, chatMessage{
		uid:        uid,
		collection: uuid.MustParse(prompt.Collection),
		prompt:     prompt.Prompt,
		completion: completion.Text,
		pageIDs:    bg.pageIDs,
	})

	return completion, nil
}
