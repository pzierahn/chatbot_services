package table

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/chatbot_services/auth"
	"github.com/pzierahn/chatbot_services/documents"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
)

type Service struct {
	pb.UnimplementedTableServiceServer
	db       *pgxpool.Pool
	auth     auth.Service
	document *documents.Service
	agent    llm.Completion
}

func New(db *pgxpool.Pool, auth auth.Service, document *documents.Service, agent llm.Completion) *Service {
	return &Service{
		db:       db,
		auth:     auth,
		document: document,
		agent:    agent,
	}
}
