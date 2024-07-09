package collections

import (
	"cloud.google.com/go/storage"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/search"
	"github.com/pzierahn/chatbot_services/services/account"
	pb "github.com/pzierahn/chatbot_services/services/proto"
)

type Service struct {
	pb.UnimplementedCollectionServiceServer
	Auth     account.Verifier
	Database *datastore.Service
	Storage  *storage.BucketHandle
	Search   search.Index
}
