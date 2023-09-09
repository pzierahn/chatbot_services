package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/braingain"
	"github.com/pzierahn/braingain/database"
	pb "github.com/pzierahn/braingain/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"sort"
	"strings"
)

type Server struct {
	pb.UnimplementedBraingainServer
	db   *database.Client
	chat *braingain.Chat
}

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

func (server *Server) ListDocuments(ctx context.Context, _ *emptypb.Empty) (*pb.Documents, error) {

	docs, err := server.db.ListDocuments(ctx)
	if err != nil {
		return nil, err
	}

	var documents pb.Documents
	for _, doc := range docs {
		documents.Items = append(documents.Items, &pb.Documents_Document{
			Id:       doc.Id.String(),
			Filename: doc.Filename,
			Pages:    uint32(doc.Pages),
		})
	}

	sort.Slice(documents.Items, func(i, j int) bool {
		return documents.Items[i].Filename < documents.Items[j].Filename
	})

	return &documents, nil
}

func (server *Server) FindDocuments(ctx context.Context, req *pb.DocumentQuery) (*pb.Documents, error) {

	docs, err := server.db.FindDocuments(ctx, "%"+req.Query+"%")
	if err != nil {
		return nil, err
	}

	var documents pb.Documents
	for _, doc := range docs {
		documents.Items = append(documents.Items, &pb.Documents_Document{
			Id:       doc.Id.String(),
			Filename: doc.Filename,
			Pages:    uint32(doc.Pages),
		})
	}

	sort.Slice(documents.Items, func(i, j int) bool {
		return documents.Items[i].Filename < documents.Items[j].Filename
	})

	return &documents, nil
}

func NewServer(db *database.Client, chat *braingain.Chat) *Server {
	return &Server{
		chat: chat,
		db:   db,
	}
}
