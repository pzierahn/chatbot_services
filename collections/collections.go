package collections

import (
	"cloud.google.com/go/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/auth"
	pb "github.com/pzierahn/brainboost/proto"
)

const (
	bucket = "documents"
)

type Service struct {
	pb.UnimplementedCollectionServiceServer
	auth    auth.Service
	db      *pgxpool.Pool
	storage *storage.BucketHandle
}

func NewServer(auth auth.Service, db *pgxpool.Pool, storage *storage.BucketHandle) *Service {
	return &Service{
		db:      db,
		storage: storage,
		auth:    auth,
	}
}
