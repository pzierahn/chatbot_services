package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/auth"
	"github.com/pzierahn/braingain/braingain"
	"github.com/pzierahn/braingain/database"
	pb "github.com/pzierahn/braingain/proto"
	"log"
	"sort"
	"strings"
)

type background struct {
	text []string
	docs []*pb.Completion_Document
}

func (server *Server) getBackgroundFromPrompt(ctx context.Context, uid uuid.UUID, prompt *pb.Prompt) (*background, error) {
	sort.Slice(prompt.Documents, func(i, j int) bool {
		return prompt.Documents[i].Filename < prompt.Documents[j].Filename
	})

	for _, doc := range prompt.Documents {
		sort.Slice(doc.Pages, func(i, j int) bool {
			return doc.Pages[i] < doc.Pages[j]
		})
	}

	var text []string
	var sources []*pb.Completion_Document

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
		}

		text = append(text, strings.Join(parts, "\n"))
		sources = append(sources, &pb.Completion_Document{
			Id:       doc.Id,
			Filename: doc.Filename,
			Pages:    doc.Pages,
		})
	}

	return &background{
		text: text,
		docs: sources,
	}, nil
}

func (server *Server) getBackgroundFromDB(ctx context.Context, uid uuid.UUID, prompt *pb.Prompt) (*background, error) {

	collection, err := uuid.Parse(prompt.Collection)
	if err != nil {
		return nil, err
	}

	query := braingain.SearchQuery{
		UserId:     uid.String(),
		Collection: &collection,
		Prompt:     prompt.Prompt,
		Limit:      int(prompt.Options.Limit),
		Threshold:  prompt.Options.Threshold,
	}

	results, err := server.chat.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	var sources []*pb.Completion_Document
	var fragments []string

	for _, doc := range results {
		text := make([]string, len(doc.Pages))
		pages := make([]uint32, len(doc.Pages))
		scores := make([]float32, len(doc.Pages))

		for iny, page := range doc.Pages {
			text[iny] = page.Text
			pages[iny] = page.Page
			scores[iny] = page.Score
		}

		fragments = append(fragments, strings.Join(text, "\n"))
		sources = append(sources, &pb.Completion_Document{
			Id:       doc.DocId.String(),
			Filename: doc.Filename,
			Pages:    pages,
			Scores:   scores,
		})
	}

	return &background{
		text: fragments,
		docs: sources,
	}, nil
}

func (server *Server) Chat(ctx context.Context, prompt *pb.Prompt) (*pb.Completion, error) {
	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	byt, _ := json.MarshalIndent(prompt, "", "  ")

	log.Printf("Chat: %s", byt)

	if prompt.Options == nil {
		return nil, fmt.Errorf("options missing")
	}

	var bg *background

	if prompt.Documents == nil || len(prompt.Documents) == 0 {
		bg, err = server.getBackgroundFromDB(ctx, *uid, prompt)
	} else {
		bg, err = server.getBackgroundFromPrompt(ctx, *uid, prompt)
	}

	if err != nil {
		return nil, err
	}

	message := braingain.Prompt{
		Prompt:      prompt.Prompt,
		Model:       prompt.Options.Model,
		Temperature: prompt.Options.Temperature,
		MaxTokens:   int(prompt.Options.MaxTokens),
		Background:  bg.text,
	}

	response, err := server.chat.Chat(ctx, message)
	if err != nil {
		return nil, err
	}

	completion := &pb.Completion{
		Prompt:    prompt,
		Text:      response.Completion,
		Documents: bg.docs,
	}

	return completion, nil
}
