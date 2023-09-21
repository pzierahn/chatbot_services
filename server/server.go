package server

import (
	"github.com/pzierahn/brainboost/database"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
	storage_go "github.com/supabase-community/storage-go"
)

type Server struct {
	pb.UnimplementedBraingainServer
	db      *database.Client
	gpt     *openai.Client
	storage *storage_go.Client
}

func NewServer(db *database.Client, gpt *openai.Client, storage *storage_go.Client) *Server {
	return &Server{
		gpt:     gpt,
		db:      db,
		storage: storage,
	}
}
