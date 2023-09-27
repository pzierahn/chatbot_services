package collections

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/auth"
	pb "github.com/pzierahn/brainboost/proto"
	storage_go "github.com/supabase-community/storage-go"
)

const (
	bucket = "documents"
)

type Service struct {
	pb.UnimplementedCollectionServiceServer
	auth    auth.Service
	db      *pgxpool.Pool
	storage *storage_go.Client
}

func NewServer(auth auth.Service, db *pgxpool.Pool, storage *storage_go.Client) *Service {
	return &Service{
		db:      db,
		storage: storage,
		auth:    auth,
	}
}
