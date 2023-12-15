package chat

import (
	"context"
	pb "github.com/pzierahn/brainboost/proto"
	"sort"
	"strings"
)

func (service *Service) searchForContext(ctx context.Context, prompt *pb.Prompt) ([]string, []string, error) {

	results, err := service.docs.Search(ctx, &pb.SearchQuery{
		CollectionId: prompt.CollectionId,
		Query:        prompt.Prompt,
		Limit:        prompt.Limit,
		Threshold:    prompt.Threshold,
	})
	if err != nil {
		return nil, nil, err
	}

	sort.Slice(results.Items, func(i, j int) bool {
		return results.Items[i].Page > results.Items[j].Page
	})

	// Map documentIds to Content
	text := make(map[string][]string)
	chunkIds := make([]string, len(results.Items))

	for inx, chunk := range results.Items {
		docId := chunk.DocumentId
		text[docId] = append(text[docId], chunk.Content)
		chunkIds[inx] = chunk.Id
	}

	contextTexts := make([]string, len(text))
	var inx int
	for _, docId := range text {
		contextTexts[inx] = strings.Join(docId, "\n")
		inx++
	}

	return chunkIds, contextTexts, nil
}
