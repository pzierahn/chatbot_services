package documents

import (
	"cloud.google.com/go/storage"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/vectordb"
)

type Service struct {
	pb.UnimplementedDocumentServiceServer
	Auth        account.Verifier
	Database    *datastore.Service
	Storage     *storage.BucketHandle
	SearchIndex vectordb.DB
}
