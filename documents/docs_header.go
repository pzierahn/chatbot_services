package documents

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (service *Service) GetHeader(ctx context.Context, req *pb.DocumentID) (*pb.DocumentHeader, error) {
	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	doc := &pb.DocumentHeader{
		Id: req.Id,
	}

	var timestamp time.Time
	var meta DocumentMeta

	err = service.db.QueryRow(ctx,
		`SELECT collection_id, created_at, metadata
				FROM documents
				WHERE id = $1 AND 
				      user_id = $2`,
		req.Id, userId,
	).Scan(
		&doc.CollectionId,
		&timestamp,
		&meta,
	)
	if err != nil {
		return nil, err
	}

	doc.Metadata = meta.toProto()
	doc.CreatedAt = timestamppb.New(timestamp)

	return doc, nil
}
