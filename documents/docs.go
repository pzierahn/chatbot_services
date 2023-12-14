package documents

import (
	"cloud.google.com/go/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/auth"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/pzierahn/brainboost/vectordb"
	"github.com/sashabaranov/go-openai"
)

const (
	bucket          = "documents"
	embeddingsModel = openai.AdaEmbeddingV2
)

type Service struct {
	pb.UnimplementedDocumentServiceServer
	auth     auth.Service
	account  *account.Service
	db       *pgxpool.Pool
	gpt      *openai.Client
	storage  *storage.BucketHandle
	vectorDB *vectordb.DB
}

type Config struct {
	Auth     auth.Service
	Account  *account.Service
	DB       *pgxpool.Pool
	GPT      *openai.Client
	Storage  *storage.BucketHandle
	VectorDB *vectordb.DB
}

func FromConfig(config *Config) *Service {
	return &Service{
		auth:     config.Auth,
		gpt:      config.GPT,
		db:       config.DB,
		storage:  config.Storage,
		account:  config.Account,
		vectorDB: config.VectorDB,
	}
}
