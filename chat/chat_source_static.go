package chat

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/utils"
	"sort"
)

type documentPages struct {
	documentIds  []string
	collectionId string
	userId       string
}

func (service *Service) getDocumentsChunkIds(ctx context.Context, query documentPages) (*pb.ReferenceIDs, error) {
	rows, err := service.db.Query(ctx,
		`SELECT chunk.id
		FROM document_chunks as chunk, documents as doc
		WHERE
		    document_id = ANY($1) AND
		    doc.id = chunk.document_id AND
		    doc.collection_id = $2 AND
		    user_id = $3
		ORDER BY index`,
		query.documentIds, query.collectionId, query.userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ids := &pb.ReferenceIDs{}
	for rows.Next() {
		var id string

		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		ids.Items = append(ids.Items, id)
	}

	return ids, nil
}

func (service *Service) getDocumentsContext(ctx context.Context, userId string, prompt *pb.ThreadPrompt) (*chunks, error) {
	sort.Slice(prompt.DocumentIds, func(i, j int) bool {
		return prompt.DocumentIds[i] < prompt.DocumentIds[j]
	})

	refIds, err := service.getDocumentsChunkIds(ctx, documentPages{
		documentIds:  prompt.DocumentIds,
		collectionId: prompt.CollectionId,
		userId:       userId,
	})
	if err != nil {
		return nil, err
	}

	refs, err := service.docs.GetReferences(ctx, refIds)
	if err != nil {
		return nil, err
	}

	data := &chunks{}

	for _, ref := range refs.Items {
		for _, chunk := range ref.Chunks {
			data.ids = append(data.ids, chunk.Id)
			data.texts = append(data.texts, chunk.Text)

			source := fmt.Sprintf("%s p.%d", utils.GetDocumentTitle(ref.Metadata), chunk.Index+1)
			data.source = append(data.source, source)
		}
	}

	return data, nil
}
