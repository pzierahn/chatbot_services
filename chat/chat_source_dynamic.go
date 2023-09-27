package chat

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/pzierahn/brainboost/proto"
	"sort"
	"strings"
)

func (service *Service) getSourceFromDB(ctx context.Context, prompt *pb.Prompt) (*chatContext, error) {

	query := &pb.SearchQuery{
		CollectionId: prompt.CollectionId,
		Query:        prompt.Prompt,
		Limit:        prompt.Options.Limit,
		Threshold:    prompt.Options.Threshold,
	}

	results, err := service.docs.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	sort.Slice(results.Items, func(i, j int) bool {
		return results.Items[i].Page > results.Items[j].Page
	})

	bg := chatContext{}

	filename := make(map[string]string)
	pageIds := make(map[string][]uuid.UUID)
	pages := make(map[string][]uint32)
	scores := make(map[string][]float32)
	text := make(map[string][]string)

	for _, doc := range results.Items {
		docId := doc.DocumentId

		filename[docId] = doc.Filename
		pages[docId] = append(pages[docId], doc.Page)
		scores[docId] = append(scores[docId], doc.Score)
		text[docId] = append(text[docId], doc.Content)
		pageIds[docId] = append(pageIds[docId], uuid.MustParse(doc.Id))
	}

	for docId := range filename {
		bg.fragments = append(bg.fragments, strings.Join(text[docId], "\n"))
		bg.docs = append(bg.docs, &pb.ChatMessage_Document{
			Id:       docId,
			Filename: filename[docId],
			Pages:    pages[docId],
			Scores:   scores[docId],
		})
		bg.pageIds = pageIds[docId]
	}

	return &bg, nil
}
