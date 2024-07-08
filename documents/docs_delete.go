package documents

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/proto"
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

	var ids []string
	for _, chunk := range doc.Content {
		ids = append(ids, chunk.Id.String())
	}

	err = service.SearchIndex.Delete(ids)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
