package documents

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/proto"
)

func (service *Service) List(ctx context.Context, req *pb.DocumentFilter) (*pb.DocumentList, error) {

	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	filter := datastore.DocumentFilter{
		UserId: userId,
		Query:  req.Query,
	}

	filter.CollectionId, err = uuid.Parse(req.CollectionId)
	if err != nil {
		return nil, err
	}

	docs, err := service.Database.ListDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := &pb.DocumentList{
		Items: make(map[string]*pb.DocumentMetadata),
	}

	for _, doc := range docs {

		var metadata *pb.DocumentMetadata

		switch doc.Type {
		case datastore.DocumentTypeWeb:
			metadata = &pb.DocumentMetadata{
				Data: &pb.DocumentMetadata_Web{
					Web: &pb.Webpage{
						Title: doc.Name,
						Url:   doc.Source,
					},
				},
			}
		case datastore.DocumentTypePDF:
			metadata = &pb.DocumentMetadata{
				Data: &pb.DocumentMetadata_File{
					File: &pb.File{
						Filename: doc.Name,
						Path:     doc.Source,
					},
				},
			}
		default:
			return nil, fmt.Errorf("unknown document type: %s", doc.Type)
		}

		result.Items[doc.Id.String()] = metadata
	}

	return result, nil
}
