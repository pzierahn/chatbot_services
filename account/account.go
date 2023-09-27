package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/auth"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service interface {
	GetModelUsages(context.Context, *emptypb.Empty) (*pb.ModelUsages, error)
	CreateUsage(ctx context.Context, usage Usage) (uuid.UUID, error)
}

type ServiceImpl struct {
	pb.UnimplementedAccountServiceServer
	db   *pgxpool.Pool
	auth auth.Service
}

type Config struct {
	Auth auth.Service
	DB   *pgxpool.Pool
}

func FromConfig(config *Config) Service {
	return &ServiceImpl{
		db:   config.DB,
		auth: config.Auth,
	}
}
