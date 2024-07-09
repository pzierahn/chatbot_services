package collections

import (
	"cloud.google.com/go/storage"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/search"
)

type Service struct {
	pb.UnimplementedCollectionServiceServer
	Auth     account.Verifier
	Database *datastore.Service
	Storage  *storage.BucketHandle
	Search   search.Index
}
