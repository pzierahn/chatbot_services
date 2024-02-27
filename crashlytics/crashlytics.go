package crashlytics

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/chatbot_services/auth"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	pb.UnimplementedCrashlyticsServiceServer
	auth auth.Service
	db   *pgxpool.Pool
}

func New(auth auth.Service, db *pgxpool.Pool) *Service {
	return &Service{
		db:   db,
		auth: auth,
	}
}

func (server *Service) RecordError(ctx context.Context, req *pb.Error) (*emptypb.Empty, error) {
	uid, err := server.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	_, err = server.db.Exec(
		ctx,
		`insert into crashlytics (user_id, exception, stack_trace, app_version)
			values ($1, $2, $3, $4)
			returning id`,
		uid,
		req.Exception,
		req.StackTrace,
		req.AppVersion)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
