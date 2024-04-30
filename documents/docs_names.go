package documents

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func (service *Service) MapDocumentNames(ctx context.Context, collection *pb.CollectionID) (*pb.DocumentNames, error) {
	documentList, err := service.List(ctx, &pb.DocumentFilter{
		CollectionId: collection.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list documents: %v", err)
	}

	documentNames := make(map[string]string)

	for id, document := range documentList.Items {
		switch document.Data.(type) {
		case *pb.DocumentMetadata_Web:
			documentNames[document.GetWeb().Title] = id
		case *pb.DocumentMetadata_File:
			filename := strings.TrimSuffix(document.GetFile().Filename, ".pdf")
			documentNames[filename] = id
		}
	}

	return &pb.DocumentNames{
		Items: documentNames,
	}, nil
}
