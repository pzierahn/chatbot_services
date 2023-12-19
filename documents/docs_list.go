package documents

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"sort"
)

func (service *Service) List(ctx context.Context, req *pb.DocumentFilter) (*pb.Documents, error) {

	userId, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := service.db.Query(ctx,
		`SELECT document_id, filename, max(page)
		FROM documents AS doc
		    join document_chunks AS em on doc.id = em.document_id
		WHERE
		    doc.user_id = $1 AND
		    doc.collection_id = $2::uuid AND
		    doc.filename LIKE $3
		GROUP BY document_id, filename, collection_id`,
		userId, req.CollectionId, "%"+req.Query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents pb.Documents
	for rows.Next() {
		source := pb.Documents_Document{}

		err = rows.Scan(
			&source.Id,
			&source.Filename,
			&source.Pages)
		if err != nil {
			return nil, err
		}

		documents.Items = append(documents.Items, &source)
	}

	sort.Slice(documents.Items, func(i, j int) bool {
		return documents.Items[i].Filename < documents.Items[j].Filename
	})

	return &documents, nil
}
