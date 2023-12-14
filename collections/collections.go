package collections

import (
	"cloud.google.com/go/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/auth"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/pzierahn/brainboost/vectordb"
)

const (
	bucket = "documents"
)

type Service struct {
	pb.UnimplementedCollectionServiceServer
	auth     auth.Service
	db       *pgxpool.Pool
	storage  *storage.BucketHandle
	vectorDB *vectordb.DB
}

func NewServer(auth auth.Service, db *pgxpool.Pool, storage *storage.BucketHandle, vectorDB *vectordb.DB) *Service {
	return &Service{
		db:       db,
		storage:  storage,
		auth:     auth,
		vectorDB: vectorDB,
	}
}
