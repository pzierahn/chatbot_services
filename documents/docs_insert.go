package documents

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
)

type document struct {
	userId    string
	document  *pb.IndexJob
	chunkMeta []*pb.Chunk
}

func (service *Service) insertIntoDB(ctx context.Context, data document) error {
	tx, err := service.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO documents (id, user_id, collection_id, metadata) VALUES ($1, $2, $3, $4)`,
		data.document.Id, data.userId, data.document.CollectionId, data.document.Document)
	if err != nil {
		return err
	}

	for inx := 0; inx < len(data.chunkMeta); inx++ {
		chunk := data.chunkMeta[inx]
		_, err = tx.Exec(
			ctx,
			`INSERT INTO document_chunks (id, document_id, text, metadata) VALUES ($1, $2, $3, $4)`,
			chunk.Id, data.document.Id, chunk.Text, chunk.Metadata)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
