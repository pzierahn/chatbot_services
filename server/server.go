package server

import (
	"github.com/pzierahn/braingain/database"
	"github.com/pzierahn/braingain/index"
	pb "github.com/pzierahn/braingain/proto"
	"github.com/sashabaranov/go-openai"
	storage_go "github.com/supabase-community/storage-go"
)

type Server struct {
	pb.UnimplementedBraingainServer
	db      *database.Client
	gpt     *openai.Client
	storage *storage_go.Client
	index   index.Index
}

func NewServer(db *database.Client, gpt *openai.Client, storage *storage_go.Client) *Server {
	return &Server{
		gpt:     gpt,
		db:      db,
		storage: storage,
		index: index.Index{
			DB:      db,
			GPT:     gpt,
			Storage: storage,
		},
	}
}
