package server

import (
	"context"
	pb "github.com/pzierahn/braingain/proto"
	"sort"
)

func (server *Server) FindDocuments(ctx context.Context, req *pb.DocumentQuery) (*pb.Documents, error) {

	docs, err := server.db.FindDocuments(ctx, patrick.String(), "%"+req.Query+"%")
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
