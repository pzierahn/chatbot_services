package collections

import (
	"cloud.google.com/go/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"github.com/pzierahn/brainboost/auth"
	pb "github.com/pzierahn/brainboost/proto"
)

const (
	bucket = "documents"
)

type Service struct {
	pb.UnimplementedCollectionServiceServer
	auth     auth.Service
	db       *pgxpool.Pool
	storage  *storage.BucketHandle
	pinecone pinecone_grpc.VectorServiceClient
}

func NewServer(auth auth.Service, db *pgxpool.Pool, storage *storage.BucketHandle, pinecone pinecone_grpc.VectorServiceClient) *Service {
	return &Service{
		db:       db,
		storage:  storage,
		auth:     auth,
		pinecone: pinecone,
	}
}
