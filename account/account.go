package account

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/chatbot_services/auth"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/grpc/status"
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

// NoFundingCode is the error code returned when a user has no founding. https://grpc.github.io/grpc/core/md_doc_statuscodes.html
const NoFundingCode = 17

func NoFundingError() error {
	return status.Errorf(NoFundingCode, "no funding")
}
