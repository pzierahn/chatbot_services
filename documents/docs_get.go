package documents

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (service *Service) Get(ctx context.Context, req *pb.DocumentID) (*pb.Document, error) {
	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := service.db.Query(ctx,
		`SELECT doc.id, chunk.id, chunk.text, chunk.index
				FROM document_chunks as chunk,
				     documents as doc
				WHERE doc.id = $1 AND
				      chunk.document_id = doc.id AND
				      doc.user_id = $2
			  	ORDER BY index`,
		req.Id, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []*pb.Chunk
	for rows.Next() {
		var docId string
		var chunk pb.Chunk

		err = rows.Scan(
			&docId,
			&chunk.Id,
			&chunk.Text,
			&chunk.Index,
		)
		if err != nil {
			return nil, err
		}

		chunks = append(chunks, &chunk)
	}

	doc := &pb.Document{
		Id:     req.Id,
		Chunks: chunks,
	}
	var timestamp time.Time
	var meta DocumentMeta

	err = service.db.QueryRow(ctx,
		`SELECT collection_id, created_at, metadata
				FROM documents
				WHERE id = $1`, req.Id).Scan(
		&doc.CollectionId,
		&timestamp,
		&meta,
	)
	if err != nil {
		return nil, err
	}

	doc.Metadata = metaToProto(meta)
	doc.CreatedAt = timestamppb.New(timestamp)

	return doc, nil
}
