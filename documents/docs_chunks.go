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

	rows, err := service.db.Query(ctx,
		`SELECT document_id, chunks.id, index
				FROM document_chunks as chunks,
				     documents as docs
				WHERE chunks.id = ANY($1) AND
				      chunks.document_id = docs.id AND
				      docs.user_id = $2
			  	ORDER BY index`,
		req.Items, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chunksMapping := make(map[string][]*pb.Chunk)
	for rows.Next() {
		var docId string
		var chunk pb.Chunk

		err = rows.Scan(&docId, &chunk.Id, &chunk.Index)
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
		var meta DocumentMeta

		err := service.db.QueryRow(ctx,
			`SELECT collection_id, created_at, metadata
				FROM documents
				WHERE id = $1`, docId).Scan(
			&doc.CollectionId,
			&timestamp,
			&meta,
		)
		if err != nil {
			return nil, err
		}

		doc.Metadata = metaToProto(meta)
		doc.CreatedAt = timestamppb.New(timestamp)

		references.Items = append(references.Items, doc)
	}

	return &references, nil
}
