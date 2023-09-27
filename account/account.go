package account

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/auth"
	pb "github.com/pzierahn/brainboost/proto"
)

type Service struct {
	pb.UnimplementedAccountServiceServer
	db   *pgxpool.Pool
	auth auth.Service
}

type Config struct {
	Auth auth.Service
	DB   *pgxpool.Pool
}

func FromConfig(config *Config) *Service {
	return &Service{
		db:   config.DB,
		auth: config.Auth,
	}
}
