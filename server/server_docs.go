package server

import (
	"context"
	pb "github.com/pzierahn/braingain/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"sort"
)

func (server *Server) ListDocuments(ctx context.Context, _ *emptypb.Empty) (*pb.Documents, error) {

	docs, err := server.db.ListDocuments(ctx)
	if err != nil {
		return nil, err
	}

	var documents pb.Documents
	for _, doc := range docs {
		documents.Items = append(documents.Items, &pb.Documents_Document{
			Id:       doc.Id.String(),
			Filename: doc.Filename,
			Pages:    uint32(doc.Pages),
		})
	}

	sort.Slice(documents.Items, func(i, j int) bool {
		return documents.Items[i].Filename < documents.Items[j].Filename
	})

	return &documents, nil
}

func (server *Server) FindDocuments(ctx context.Context, req *pb.DocumentQuery) (*pb.Documents, error) {

	docs, err := server.db.FindDocuments(ctx, "%"+req.Query+"%")
	if err != nil {
		return nil, err
	}

	var documents pb.Documents
	for _, doc := range docs {
		documents.Items = append(documents.Items, &pb.Documents_Document{
			Id:       doc.Id.String(),
			Filename: doc.Filename,
			Pages:    uint32(doc.Pages),
		})
	}

	sort.Slice(documents.Items, func(i, j int) bool {
		return documents.Items[i].Filename < documents.Items[j].Filename
	})

	return &documents, nil
}
