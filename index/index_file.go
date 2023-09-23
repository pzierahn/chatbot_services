package index

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/brainboost/database"
	"github.com/pzierahn/brainboost/pdf"
)

type Progress struct {
	TotalPages   int
	FinishedPage int
}

func (index Index) Process(ctx context.Context, doc DocumentId, data []byte, ch ...chan<- Progress) (uuid.UUID, error) {

	pages, err := pdf.GetPagesFromBytes(ctx, data)
	if err != nil {
		return uuid.Nil, err
	}

	embeddings, err := index.GetPagesWithEmbeddings(ctx, pages, ch...)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := index.DB.UpsertDocument(ctx, database.Document{
		UserID:       doc.UserID,
		CollectionID: doc.CollectionID,
		Filename:     doc.Filename,
		Path:         doc.path(),
		Pages:        embeddings,
	})
	if err != nil {
		return uuid.Nil, err
	}

	doc.DocumentID = id

	err = index.Upload(doc, data)
	if err != nil {
		// Delete document from database if upload fails.
		_ = index.DB.DeleteDocument(ctx, id, doc.UserID)
		return uuid.Nil, err
	}

	return id, nil
}
