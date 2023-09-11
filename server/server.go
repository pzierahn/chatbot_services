package server

import (
	"github.com/pzierahn/braingain/braingain"
	"github.com/pzierahn/braingain/database"
	pb "github.com/pzierahn/braingain/proto"
)

type Server struct {
	pb.UnimplementedBraingainServer
	db   *database.Client
	chat *braingain.Chat
}

func NewServer(db *database.Client, chat *braingain.Chat) *Server {
	return &Server{
		chat: chat,
		db:   db,
	}
}
