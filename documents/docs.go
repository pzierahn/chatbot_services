package documents

import (
	"cloud.google.com/go/storage"
	"github.com/pzierahn/chatbot_services/auth"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/vectordb"
)

type Service struct {
	pb.UnimplementedDocumentServiceServer
	Auth        auth.Service
	Database    *datastore.Service
	Storage     *storage.BucketHandle
	SearchIndex vectordb.DB
}
