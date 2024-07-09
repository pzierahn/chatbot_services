package crashlytics

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/auth"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Service implements the crashlytics service.
type Service struct {
	pb.UnimplementedCrashlyticsServiceServer
	Auth     auth.Service
	Database *datastore.Service
}

// RecordError records a frontend error in the database.
func (server *Service) RecordError(ctx context.Context, req *pb.Error) (*emptypb.Empty, error) {
	userId, err := server.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	err = server.Database.InsertError(ctx, &datastore.Error{
		Id:         uuid.New(),
		UserId:     userId,
		Exception:  req.Exception,
		StackTrace: req.StackTrace,
		AppVersion: req.AppVersion,
	})

	return &emptypb.Empty{}, nil
}
