package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/braingain"
	"github.com/pzierahn/braingain/database"
	pb "github.com/pzierahn/braingain/proto"
	"log"
	"sort"
)

type Server struct {
	pb.UnimplementedBraingainServer
	db   *database.Client
	chat *braingain.Chat
}

func (server *Server) Chat(ctx context.Context, prompt *pb.Prompt) (*pb.Completion, error) {
	log.Printf("Chat: %s", prompt.Prompt)

	response, err := server.chat.RAG(ctx, prompt.Prompt)
	if err != nil {
		return nil, err
	}

	completion := &pb.Completion{
		Prompt: prompt,
		Text:   response.Completion,
	}

	docs := make(map[uuid.UUID]*pb.Completion_Document)
	for _, doc := range response.Documents {
		if _, ok := docs[doc.Id]; ok {
			docs[doc.Id] = &pb.Completion_Document{
				Id:       doc.Id.String(),
				Filename: doc.Id.String(),
			}
		}

		docs[doc.Id].Pages = append(docs[doc.Id].Pages, uint32(doc.Page))
		docs[doc.Id].Scores = append(docs[doc.Id].Scores, doc.Score)
	}

	for _, doc := range docs {
		completion.Documents = append(completion.Documents, doc)
	}

	sort.Slice(completion.Documents, func(i, j int) bool {
		return completion.Documents[i].Filename > completion.Documents[j].Filename
	})

	return completion, nil
}

func NewServer(db *database.Client, chat *braingain.Chat) *Server {
	return &Server{
		chat: chat,
		db:   db,
	}
}
