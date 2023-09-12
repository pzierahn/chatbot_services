package index

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/database"
	"github.com/pzierahn/braingain/pdf"
)

type Progress struct {
	TotalPages   int
	FinishedPage int
}

func (index Index) Process(ctx context.Context, doc DocumentId, data []byte, ch ...chan<- Progress) (*uuid.UUID, error) {

	pages, err := pdf.GetPagesFromBytes(ctx, data)
	if err != nil {
		return nil, err
	}

	embeddings, err := index.GetPagesWithEmbeddings(ctx, pages, ch...)
	if err != nil {
		return nil, err
	}

	id, err := index.DB.UpsertDocument(ctx, database.Document{
		UserId:     doc.UserId,
		Collection: doc.Collection,
		Filename:   doc.Filename,
		Path:       doc.path(),
		Pages:      embeddings,
	})
	if err != nil {
		return nil, err
	}

	doc.DocId = *id

	err = index.Upload(doc, data)
	if err != nil {
		// Delete document from database if upload fails.
		_ = index.DB.DeleteDocument(ctx, *id, doc.UserId)
		return nil, err
	}

	return id, nil
}
