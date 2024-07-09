package account

import (
	"context"
	"github.com/pzierahn/chatbot_services/auth"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/services/proto"
)

// Verifier is an interface for verifying user identity and funding.
type Verifier interface {
	// Verify check if the context contains valid user credentials. If successful, the user ID is returned.
	Verify(context.Context) (userId string, err error)
	// VerifyFunding checks if the context contains valid user credentials and the user has enough funding to perform an action.
	VerifyFunding(context.Context) (userId string, err error)
}

// Service is the account service implementation.
type Service struct {
	pb.UnimplementedAccountServer
	Database *datastore.Service
	Auth     auth.Service
}
