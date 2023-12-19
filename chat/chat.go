package chat

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/auth"
	"github.com/pzierahn/chatbot_services/documents"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
)

type Service struct {
	pb.UnimplementedChatServiceServer
	db      *pgxpool.Pool
	openai  llm.Completion
	vertex  llm.Completion
	docs    *documents.Service
	account account.Service
	auth    auth.Service
}

type Config struct {
	DB              *pgxpool.Pool
	Openai          llm.Completion
	Vertex          llm.Completion
	DocumentService *documents.Service
	AccountService  *account.Service
	AuthService     auth.Service
}

func FromConfig(config *Config) *Service {
	return &Service{
		openai:  config.Openai,
		vertex:  config.Vertex,
		db:      config.DB,
		docs:    config.DocumentService,
		account: *config.AccountService,
		auth:    config.AuthService,
	}
}
