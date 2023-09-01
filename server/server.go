package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/braingain"
	"github.com/pzierahn/braingain/database"
	pb "github.com/pzierahn/braingain/proto"
	"log"
)

type Server struct {
	pb.UnimplementedBraingainServer
	db   *database.Client
	chat *braingain.Chat
}

func (server *Server) Chat(ctx context.Context, prompt *pb.Prompt) (*pb.ChatCompletion, error) {
	log.Printf("Chat: %s", prompt.Prompt)

	response, err := server.chat.RAG(ctx, prompt.Prompt)
	if err != nil {
		return nil, err
	}

	completion := &pb.ChatCompletion{
		Prompt:     prompt.Prompt,
		Completion: response.Completion,
	}

	for _, doc := range response.Documents {
		completion.Sources = append(completion.Sources, &pb.Source{
			Id:    doc.Source.String(),
			Score: doc.Score,
			Page:  int32(doc.Page),
		})
	}

	return completion, nil
}

func (server *Server) GetDocument(ctx context.Context, req *pb.DocumentId) (*pb.Document, error) {
	log.Printf("GetDocument: %s", req.Id)

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	doc, err := server.db.GetDocument(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.Document{
		Id:       doc.Id.String(),
		Filename: doc.Filename,
	}, nil
}

func NewServer(db *database.Client, chat *braingain.Chat) *Server {
	return &Server{
		chat: chat,
		db:   db,
	}
}
