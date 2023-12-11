package collections

import (
	"context"
)

func (server *Service) collectionDocumentIds(ctx context.Context, userId, collectionId string) ([]string, error) {

	rows, err := server.db.Query(ctx,
		`SELECT id
		FROM documents
		WHERE
		    user_id = $1 AND
		    collection_id = $2::uuid`,
		userId, collectionId)
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

func (server *Service) documentChunkIds(ctx context.Context, userId, collectionId string) ([]string, error) {

	rows, err := server.db.Query(ctx,
		`SELECT id
		FROM documents
		WHERE
		    user_id = $1 AND
		    collection_id = $2::uuid`,
		userId, collectionId)
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
