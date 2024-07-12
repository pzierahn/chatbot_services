package documents

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (service *Service) Delete(ctx context.Context, req *pb.DocumentID) (*emptypb.Empty, error) {
	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	docId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	doc, err := service.Database.GetDocument(ctx, userId, docId)
	if err != nil {
		return nil, err
	}

	if doc.Type == datastore.DocumentTypePDF {
		obj := service.Storage.Object(doc.Source)
		err = obj.Delete(ctx)
		if err != nil {
			return nil, err
		}
	}

	err = service.SearchIndex.DeleteDocument(ctx, userId, req.Id)
	if err != nil {
		return nil, err
	}

	err = service.Database.DeleteDocument(ctx, userId, docId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
