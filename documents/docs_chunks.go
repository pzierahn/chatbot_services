package documents

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
)

func (service *Service) getChunkIds(ctx context.Context, documentId string) ([]string, error) {

	rows, err := service.db.Query(ctx,
		`SELECT id FROM document_chunks WHERE document_id = $1`, documentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string

	for rows.Next() {
		var id string

		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (service *Service) GetReferences(ctx context.Context, req *pb.ReferenceIDs) (*pb.References, error) {

	var chunks pb.References

	for _, id := range req.Items {
		var chunk pb.Reference

		err := service.db.QueryRow(ctx,
			`SELECT doc.filename, chunk.id, chunk.document_id, chunk.page
				FROM document_chunks as chunk, documents as doc
				WHERE chunk.id = $1 AND chunk.document_id = doc.id LIMIT 1`,
			id).Scan(
			&chunk.Filename,
			&chunk.Id,
			&chunk.DocumentId,
			&chunk.Page,
		)
		if err != nil {
			return nil, err
		}

		chunks.Items = append(chunks.Items, &chunk)
	}

	return &chunks, nil
}
