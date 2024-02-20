package documents

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (service *Service) Delete(ctx context.Context, req *pb.DocumentID) (*emptypb.Empty, error) {
	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	doc, err := service.GetDocument(ctx, req)
	if err != nil {
		return nil, err
	}

	var meta DocumentMeta
	err = service.db.QueryRow(ctx,
		`DELETE FROM documents
       		  WHERE id = $1 AND
					user_id = $2
			  RETURNING metadata`,
		doc.Id, userId).Scan(&meta)
	if err != nil {
		return nil, err
	}

	if meta.IsFile() {
		obj := service.storage.Object(meta.File.Path)
		err = obj.Delete(ctx)
		if err != nil {
			return nil, err
		}
	}

	var ids []string
	for _, chunk := range doc.Chunks {
		ids = append(ids, chunk.Id)
	}

	err = service.vectorDB.Delete(ids)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
