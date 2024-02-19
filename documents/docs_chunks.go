package documents

import (
	"context"
	"encoding/json"
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

func parseDocumentMetadata(meta []byte) (*pb.DocumentMetadata, error) {
	documentMetadata := pb.DocumentMetadata{}
	documentMetadata.Data = &pb.DocumentMetadata_Web{}

	err := json.Unmarshal(meta, &documentMetadata)
	if err == nil {
		return &documentMetadata, nil
	}

	documentMetadata.Data = &pb.DocumentMetadata_File{}
	err = json.Unmarshal(meta, &documentMetadata)
	if err == nil {
		return &documentMetadata, nil
	}

	return nil, err
}

func (service *Service) getReferences(ctx context.Context, userId string, req *pb.ReferenceIDs) (*pb.References, error) {

	chunksMapping := make(map[string][]*pb.Chunk)
	for _, id := range req.Items {
		var docId string
		var chunk pb.Chunk

		err := service.db.QueryRow(ctx,
			`SELECT document_id, chunks.id, index
				FROM document_chunks as chunks,
				     documents as docs
				WHERE chunks.id = $1 AND
				      chunks.document_id = docs.id AND
				      docs.user_id = $2`,
			id, userId).Scan(
			&docId,
			&chunk.Id,
			&chunk.Index,
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
			Metadata: &pb.DocumentMetadata{
				Data: &pb.DocumentMetadata_Web{},
			},
		}
		var timestamp time.Time
		var meta []byte

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

		doc.Metadata, err = parseDocumentMetadata(meta)
		if err != nil {
			return nil, err
		}

		doc.CreatedAt = timestamppb.New(timestamp)

		references.Items = append(references.Items, doc)
	}

	return &references, nil
}
