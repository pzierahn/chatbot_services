package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/database"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"sort"
)

func (server *Server) ListDocuments(ctx context.Context, req *pb.DocumentFilter) (*pb.Documents, error) {

	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	query := database.DocumentQuery{
		UserId: uid.String(),
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

func (server *Server) DeleteDocument(ctx context.Context, req *pb.Document) (*emptypb.Empty, error) {
	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	err = server.db.DeleteDocument(ctx, id, uid.String())
	if err != nil {
		return nil, err
	}

	resp := server.storage.RemoveFile(bucket, []string{req.Path})
	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	return &emptypb.Empty{}, nil
}

func (server *Server) UpdateDocument(ctx context.Context, req *pb.Document) (*emptypb.Empty, error) {
	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("UpdateDocument: %+v", req)

	docID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	collection, err := uuid.Parse(req.CollectionId)
	if err != nil {
		return nil, err
	}

	err = server.db.UpdateDocumentName(ctx, database.Document{
		Id:         docID,
		UserId:     uid.String(),
		Collection: collection,
		Filename:   req.Filename,
		Path:       req.Path,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
