package server

import (
	"context"
	"github.com/pzierahn/braingain/braingain"
	pb "github.com/pzierahn/braingain/proto"
)

type Server struct {
	pb.BraingainServer
	chat *braingain.Chat
}

func (server *Server) Chat(ctx context.Context, prompt *pb.Prompt) (*pb.ChatCompletion, error) {
	response, err := server.chat.RAG(ctx, prompt.Prompt)
	if err != nil {
		return nil, err
	}

	completion := &pb.ChatCompletion{
		Completion: response.Completion,
	}

	for _, source := range response.Sources {
		completion.Sources = append(completion.Sources, &pb.Source{
			Id:       source.Id,
			Score:    source.Score,
			Filename: source.Filename,
			Page:     int32(source.Page),
		})
	}

	return completion, nil
}

func NewServer(chat *braingain.Chat) *Server {
	return &Server{
		chat: chat,
	}
}
