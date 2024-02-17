package documents

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (service *Service) GetReferences(ctx context.Context, req *pb.ReferenceIDs) (*pb.References, error) {

	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	return service.getReferences(ctx, userId, req)
}

func (service *Service) getReferences(ctx context.Context, userId string, req *pb.ReferenceIDs) (*pb.References, error) {

	chunksMapping := make(map[string][]*pb.Chunk)
	for _, id := range req.Items {
		var docId string
		var chunk pb.Chunk

		err := service.db.QueryRow(ctx,
			`SELECT document_id, chunks.id, chunks.metadata
				FROM document_chunks as chunks,
				     documents as docs
				WHERE chunks.id = $1 AND
				      chunks.document_id = docs.id AND
				      docs.user_id = $2`,
			id, userId).Scan(
			&docId,
			&chunk.Id,
			&chunk.Metadata,
		)
		if err != nil {
			return nil, err
		}

		chunksMapping[docId] = append(chunksMapping[docId], &chunk)
	}

	var references pb.References

	for docId, chunks := range chunksMapping {
		doc := &pb.Document{
			Id:     docId,
			Chunks: chunks,
		}
		var timestamp time.Time

		err := service.db.QueryRow(ctx,
			`SELECT collection_id, created_at, metadata
				FROM documents
				WHERE id = $1`, docId).Scan(
			&doc.CollectionId,
			&timestamp,
			&doc.Metadata,
		)
		if err != nil {
			return nil, err
		}

		doc.CreatedAt = timestamppb.New(timestamp)

		references.Items = append(references.Items, doc)
	}

	return &references, nil
}
