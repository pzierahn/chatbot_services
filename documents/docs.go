package documents

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/auth"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
	storage_go "github.com/supabase-community/storage-go"
)

const (
	bucket          = "documents"
	embeddingsModel = openai.AdaEmbeddingV2
)

type Service struct {
	pb.UnimplementedDocumentServiceServer
	auth    auth.Service
	account account.ServiceImpl
	db      *pgxpool.Pool
	gpt     *openai.Client
	storage *storage_go.Client
}

type Config struct {
	Auth    auth.Service
	Account *account.ServiceImpl
	DB      *pgxpool.Pool
	GPT     *openai.Client
	Storage *storage_go.Client
}

func FromConfig(config *Config) *Service {
	return &Service{
		auth:    config.Auth,
		gpt:     config.GPT,
		db:      config.DB,
		storage: config.Storage,
		account: *config.Account,
	}
}
