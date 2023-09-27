package chat

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/documents"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
)

type Service struct {
	pb.UnimplementedChatServiceServer
	db      *pgxpool.Pool
	gpt     *openai.Client
	docs    *documents.Service
	account account.Service
	auth    auth.Service
}

type Config struct {
	DB              *pgxpool.Pool
	GPT             *openai.Client
	DocumentService *documents.Service
	AccountService  account.Service
	AuthService     auth.Service
}

func FromConfig(config *Config) *Service {
	return &Service{
		gpt:     config.GPT,
		db:      config.DB,
		docs:    config.DocumentService,
		account: config.AccountService,
		auth:    config.AuthService,
	}
}
