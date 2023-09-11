package server

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/pzierahn/braingain/proto"
	"log"
	"sort"
	"strings"
)

func (server *Server) Chat(ctx context.Context, prompt *pb.Prompt) (*pb.Completion, error) {
	log.Printf("Chat: %v", prompt)

	sort.Slice(prompt.Documents, func(i, j int) bool {
		return prompt.Documents[i].Filename < prompt.Documents[j].Filename
	})

	for _, doc := range prompt.Documents {
		sort.Slice(doc.Pages, func(i, j int) bool {
			return doc.Pages[i] < doc.Pages[j]
		})
	}

	var background []string
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

		background = append(background, strings.Join(parts, "\n"))
		sources = append(sources, &pb.Completion_Document{
			Id:       doc.Id,
			Filename: doc.Filename,
			Pages:    doc.Pages,
		})
	}

	response, err := server.chat.Chat(ctx, prompt.Prompt, background)
	if err != nil {
		return nil, err
	}

	completion := &pb.Completion{
		Prompt:    prompt,
		Text:      response.Completion,
		Documents: sources,
	}

	return completion, nil
}
