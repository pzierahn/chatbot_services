package chat

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"sort"
	"strings"
)

type documentPages struct {
	id           string
	collectionId string
	userId       string
}

type documentChunk struct {
	chunkId string
	page    uint32
	text    string
}

func (service *Service) getDocumentChunks(ctx context.Context, query documentPages) ([]documentChunk, error) {
	rows, err := service.db.Query(ctx,
		`SELECT chunk.id, chunk.index, chunk.text
		FROM document_chunks as chunk, documents as doc
		WHERE
		    document_id = $1 AND
		    doc.id = chunk.document_id AND
		    doc.collection_id = $2 AND
		    user_id = $3
		ORDER BY index`,
		query.id, query.collectionId, query.userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var docChunks []documentChunk

	for rows.Next() {
		var chunk documentChunk

		err = rows.Scan(
			&chunk.chunkId,
			&chunk.page,
			&chunk.text,
		)
		if err != nil {
			return nil, err
		}

		docChunks = append(docChunks, chunk)
	}

	return docChunks, nil
}

func (service *Service) getDocumentsContext(ctx context.Context, userId string, prompt *pb.ThreadPrompt) (*chunks, error) {
	sort.Slice(prompt.DocumentIds, func(i, j int) bool {
		return prompt.DocumentIds[i] < prompt.DocumentIds[j]
	})

	data := &chunks{}

	for _, docId := range prompt.DocumentIds {
		items, err := service.getDocumentChunks(ctx, documentPages{
			id:           docId,
			collectionId: prompt.CollectionId,
			userId:       userId,
		})
		if err != nil {
			return nil, err
		}

		var texts []string
		for _, chunk := range items {
			data.ids = append(data.ids, chunk.chunkId)
			texts = append(texts, chunk.text)
		}

		data.texts = append(data.texts, strings.Join(texts, "\n"))
	}

	return data, nil
}
