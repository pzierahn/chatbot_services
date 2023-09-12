package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/database"
	pb "github.com/pzierahn/braingain/proto"
	"log"
	"sort"
)

func (server *Server) GetDocuments(ctx context.Context, req *pb.DocumentQuery) (*pb.Documents, error) {

	query := database.DocumentQuery{
		UserId: patrick.String(),
		Query:  "%" + req.Query + "%",
	}

	if req.Collection != "" {
		collId := uuid.MustParse(req.Collection)
		query.Collection = &collId
	}

	log.Printf("GetDocuments: %+v", query)

	docs, err := server.db.FindDocuments(ctx, query)
	if err != nil {
		return nil, err
	}

	var documents pb.Documents
	for _, doc := range docs {
		documents.Items = append(documents.Items, &pb.Documents_Document{
			Id:       doc.Id.String(),
			Filename: doc.Filename,
			Pages:    doc.Pages,
		})
	}

	sort.Slice(documents.Items, func(i, j int) bool {
		return documents.Items[i].Filename < documents.Items[j].Filename
	})

	return &documents, nil
}
