package pinecone_search

import (
	"context"
	"github.com/pinecone-io/go-pinecone/pinecone"
)

// DeleteCollection deletes all vectors in a collection. Delete by filter
// is not supported by Pinecone serverless, so we need to assemble a list
// with vector ids manually and delete them one by one.
func (db *Search) deleteByPrefix(ctx context.Context, prefix string) error {

	idxConnection, err := db.getIndexConnection(ctx)
	if err != nil {
		return err
	}

	defer func() { _ = idxConnection.Close() }()

	limit := uint32(100)

	var token *string
	var ids []string

	for {
		list, err := idxConnection.ListVectors(ctx, &pinecone.ListVectorsRequest{
			Prefix:          &prefix,
			Limit:           &limit,
			PaginationToken: token,
		})
		if err != nil {
			return err
		}

		for _, id := range list.VectorIds {
			ids = append(ids, *id)
		}

		token = list.NextPaginationToken
		if token == nil {
			break
		}
	}

	if len(ids) == 0 {
		return nil
	}

	return idxConnection.DeleteVectorsById(ctx, ids)
}

func (db *Search) DeleteCollection(ctx context.Context, _, collectionId string) error {
	prefix := collectionId + "#"
	return db.deleteByPrefix(ctx, prefix)
}

func (db *Search) DeleteDocument(ctx context.Context, _, collectionId, documentId string) error {
	prefix := collectionId + "#" + documentId + "#"
	return db.deleteByPrefix(ctx, prefix)
}
