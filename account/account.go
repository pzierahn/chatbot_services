package account

import (
	"context"
	"github.com/pzierahn/chatbot_services/auth"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/proto"
)

type Verifier interface {
	Verify(context.Context) (userId string, err error)
	VerifyFunding(context.Context) (userId string, err error)
}

type Service struct {
	pb.UnimplementedAccountServiceServer
	Database *datastore.Service
	Auth     auth.Service
}
