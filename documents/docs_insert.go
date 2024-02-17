package documents

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
)

type document struct {
	userId   string
	title    string
	document *pb.IndexJob
	chunks   []*pb.Chunk
}

func (service *Service) insertIntoDB(ctx context.Context, data *document) error {
	tx, err := service.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO documents (id, user_id, collection_id, title, metadata) 
			VALUES ($1, $2, $3, $4, $5)`,
		data.document.Id, data.userId, data.document.CollectionId, data.title, data.document.Document)
	if err != nil {
		return err
	}

	for inx := 0; inx < len(data.chunks); inx++ {
		chunk := data.chunks[inx]
		_, err = tx.Exec(
			ctx,
			`INSERT INTO document_chunks (id, document_id, text, index) VALUES ($1, $2, $3, $4)`,
			chunk.Id, data.document.Id, chunk.Text, chunk.Index)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
