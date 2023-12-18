package chat

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/documents"
	"github.com/pzierahn/brainboost/llm"
	pb "github.com/pzierahn/brainboost/proto"
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
