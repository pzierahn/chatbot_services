package account

import (
	"context"
	"github.com/pzierahn/chatbot_services/auth"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/proto"
)

type Service interface {
	pb.UnimplementedAccountServiceServer
	Verify(context.Context) (userId string, err error)
	VerifyFunding(context.Context) (userId string, err error)
}

type LiveService struct {
	pb.UnimplementedAccountServiceServer
	Database *datastore.Service
	Auth     auth.Service
}
