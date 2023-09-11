package server

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/pzierahn/braingain/proto"
	"log"
	"sort"
	"strings"
)

type background struct {
	text []string
	docs []*pb.Completion_Document
}

func (server *Server) getBackgroundFromPrompt(ctx context.Context, prompt *pb.Prompt) (*background, error) {
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

		content, err := server.db.GetDocumentPages(ctx, id, doc.Pages)
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

func (server *Server) getBackgroundFromDB(ctx context.Context, prompt *pb.Prompt) (*background, error) {

	results, err := server.chat.Search(ctx, prompt.Prompt)
	if err != nil {
		return nil, err
	}

	text := make([]string, len(results))
	pages := make(map[uuid.UUID][]uint32)
	scores := make(map[uuid.UUID][]float32)

	for inx, result := range results {
		text[inx] = result.Text

		source := result.Source
		pages[source] = append(pages[source], uint32(result.Page))
		scores[source] = append(scores[source], result.Score)
	}

	var sources []*pb.Completion_Document
	for source := range pages {
		doc, err := server.db.GetDocument(ctx, source)
		if err != nil {
			return nil, err
		}

		sources = append(sources, &pb.Completion_Document{
			Id:       doc.Id.String(),
			Filename: doc.Filename,
			Pages:    pages[source],
			Scores:   scores[source],
		})
	}

	sort.Slice(sources, func(i, j int) bool {
		return sources[i].Filename < sources[j].Filename
	})

	return &background{
		text: text,
		docs: sources,
	}, nil
}

func (server *Server) Chat(ctx context.Context, prompt *pb.Prompt) (*pb.Completion, error) {
	log.Printf("Chat: %v", prompt)

	var bg *background
	var err error

	if prompt.Documents == nil || len(prompt.Documents) == 0 {
		bg, err = server.getBackgroundFromDB(ctx, prompt)
	} else {
		bg, err = server.getBackgroundFromPrompt(ctx, prompt)
	}

	if err != nil {
		return nil, err
	}

	response, err := server.chat.Chat(ctx, prompt.Prompt, bg.text)
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
