package chat

import (
	"context"
	pb "github.com/pzierahn/brainboost/proto"
	"sort"
	"strings"
)

type chunks struct {
	ids    []string
	texts  []string
	scores []float32
}

func (service *Service) searchForContext(ctx context.Context, prompt *pb.Prompt) (*chunks, error) {

	results, err := service.docs.Search(ctx, &pb.SearchQuery{
		CollectionId: prompt.CollectionId,
		Query:        prompt.Prompt,
		Limit:        prompt.Limit,
		Threshold:    prompt.Threshold,
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(results.Items, func(i, j int) bool {
		return results.Items[i].Page > results.Items[j].Page
	})

	data := &chunks{
		ids:    make([]string, len(results.Items)),
		scores: make([]float32, len(results.Items)),
	}

	// Map documentIds to Content
	text := make(map[string][]string)

	for inx, chunk := range results.Items {
		docId := chunk.DocumentId
		text[docId] = append(text[docId], chunk.Content)
		data.ids[inx] = chunk.Id
		data.scores[inx] = chunk.Score
	}

	data.texts = make([]string, len(text))
	var inx int
	for _, docId := range text {
		data.texts[inx] = strings.Join(docId, "\n")
		inx++
	}

	return data, nil
}
