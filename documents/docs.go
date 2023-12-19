package documents

import (
	"cloud.google.com/go/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/auth"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/vectordb"
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
	storage    *storage.BucketHandle
	vectorDB   *vectordb.DB
}

type Config struct {
	Auth       auth.Service
	Account    *account.Service
	DB         *pgxpool.Pool
	Embeddings llm.Embedding
	Storage    *storage.BucketHandle
	VectorDB   *vectordb.DB
}

func FromConfig(config *Config) *Service {
	return &Service{
		auth:       config.Auth,
		db:         config.DB,
		embeddings: config.Embeddings,
		storage:    config.Storage,
		account:    config.Account,
		vectorDB:   config.VectorDB,
	}
}
