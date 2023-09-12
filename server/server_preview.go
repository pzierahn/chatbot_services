package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/index"
	"github.com/pzierahn/braingain/pdf"
	pb "github.com/pzierahn/braingain/proto"
)

func (server *Server) GetDocumentPreview(ctx context.Context, ref *pb.StorageRef) (*pb.Preview, error) {
	source := index.Index{
		DB:      server.db,
		GPT:     server.gpt,
		Storage: server.storage,
	}

	byt, err := source.Download(index.DocumentId{
		UserId:     patrick.String(),
		Collection: uuid.MustParse(ref.Collection),
		DocId:      uuid.MustParse(ref.Id),
	})
	if err != nil {
		return nil, err
	}

	image, err := pdf.GetPageAsImage(ctx, byt)
	if err != nil {
		return nil, err
	}

	return &pb.Preview{
		Image: image,
	}, nil
}
