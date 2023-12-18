package documents

import (
	"cloud.google.com/go/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/llm"
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
	auth       auth.Service
	account    *account.Service
	db         *pgxpool.Pool
	embeddings llm.Embedding
	completion llm.Completion
	storage    *storage.BucketHandle
	vectorDB   *vectordb.DB
}

type Config struct {
	Auth       auth.Service
	Account    *account.Service
	DB         *pgxpool.Pool
	Embeddings llm.Embedding
	Completion llm.Completion
	Storage    *storage.BucketHandle
	VectorDB   *vectordb.DB
}

func FromConfig(config *Config) *Service {
	return &Service{
		auth:       config.Auth,
		db:         config.DB,
		embeddings: config.Embeddings,
		completion: config.Completion,
		storage:    config.Storage,
		account:    config.Account,
		vectorDB:   config.VectorDB,
	}
}
