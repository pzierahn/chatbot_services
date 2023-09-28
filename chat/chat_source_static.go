package chat

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/pzierahn/brainboost/proto"
	"sort"
	"strings"
)

type SearchQuery struct {
	UserId       uuid.UUID
	CollectionId uuid.UUID
	Prompt       string
	Limit        int
	Threshold    float32
}

type chatContext struct {
	fragments []string
	docs      []*pb.ChatMessage_Document
	pageIds   []uuid.UUID
}

type PageContentQuery struct {
	DocumentId   string
	CollectionId string
	UserId       string
	Pages        []uint32
}

func (service *Service) getPageContent(ctx context.Context, query PageContentQuery) (string, *pb.ChatMessage_Document, error) {
	rows, err := service.db.Query(ctx,
		`SELECT doc.filename, dm.page, dm.text
		FROM document_embeddings as dm, documents as doc
		WHERE
		    document_id = $1 AND
		    doc.collection_id = $2 AND
		    doc.id = dm.document_id AND
		    user_id = $3 AND
		    page = ANY($4)
		ORDER BY filename, page`,
		query.DocumentId, query.CollectionId, query.UserId, query.Pages)
	if err != nil {
		return "", nil, err
	}

	defer rows.Close()

	var (
		doc       pb.ChatMessage_Document
		fragments []string
	)

	for rows.Next() {
		var (
			page uint32
			text string
		)

		err = rows.Scan(
			&doc.Filename,
			&page,
			&text)
		if err != nil {
			return "", nil, err
		}

		doc.Pages = append(doc.Pages, page)
		fragments = append(fragments, text)
	}

	return strings.Join(fragments, "\n"), &doc, nil
}

func (service *Service) getBackgroundFromPrompt(ctx context.Context, userID uuid.UUID, prompt *pb.Prompt) (*chatContext, error) {
	sort.Slice(prompt.Documents, func(i, j int) bool {
		return prompt.Documents[i].Filename < prompt.Documents[j].Filename
	})

	for _, doc := range prompt.Documents {
		sort.Slice(doc.Pages, func(i, j int) bool {
			return doc.Pages[i] < doc.Pages[j]
		})
	}

	bg := chatContext{}
	for _, doc := range prompt.Documents {
		fragment, content, err := service.getPageContent(ctx, PageContentQuery{
			DocumentId:   doc.Id,
			CollectionId: prompt.CollectionId,
			UserId:       userID.String(),
			Pages:        doc.Pages,
		})
		if err != nil {
			return nil, err
		}

		bg.fragments = append(bg.fragments, fragment)
		bg.docs = append(bg.docs, content)
		bg.pageIds = append(bg.pageIds, uuid.MustParse(doc.Id))
	}

	return &bg, nil
}
