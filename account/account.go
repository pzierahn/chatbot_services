package account

import (
	"github.com/pzierahn/chatbot_services/auth"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedAccountServiceServer
	Database *datastore.Service
	Auth     auth.Service
}

// NoFundingCode is the error code returned when a user has no founding. https://grpc.github.io/grpc/core/md_doc_statuscodes.html
const NoFundingCode = 17

func NoFundingError() error {
	return status.Errorf(NoFundingCode, "no funding")
}
