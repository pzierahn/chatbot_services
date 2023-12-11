package documents

import (
	"context"
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
