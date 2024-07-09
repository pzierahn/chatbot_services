package collections

import (
	"cloud.google.com/go/storage"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/search"
	"github.com/pzierahn/chatbot_services/services/account"
)

type Service struct {
	pb.UnimplementedCollectionServiceServer
	Auth     account.Verifier
	Database *datastore.Service
	Storage  *storage.BucketHandle
	Search   search.Index
}
