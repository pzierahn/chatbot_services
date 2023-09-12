package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/database"
	"github.com/pzierahn/braingain/index"
	pb "github.com/pzierahn/braingain/proto"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (server *Server) DeleteDocument(ctx context.Context, req *pb.StorageRef) (*emptypb.Empty, error) {
	id := uuid.MustParse(req.Id)
	uid := patrick.String()

	err := server.db.DeleteDocument(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	source := index.Index{
		DB:      server.db,
		GPT:     server.gpt,
		Storage: server.storage,
	}

	col := uuid.MustParse(req.Collection)

	err = source.Delete(index.DocumentId{
		UserId:     patrick.String(),
		Collection: col,
		DocId:      id,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
