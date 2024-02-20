package chat

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
)

type chunks struct {
	ids    []string
	texts  []string
	scores map[string]float32
}

func (service *Service) searchForContext(ctx context.Context, prompt *pb.ThreadPrompt) (*chunks, error) {

	if prompt.Limit == 0 {
		return &chunks{}, nil
	}

	results, err := service.docs.Search(ctx, &pb.SearchQuery{
		CollectionId: prompt.CollectionId,
		Query:        prompt.Prompt,
		Limit:        prompt.Limit,
		Threshold:    prompt.Threshold,
	})
	if err != nil {
		return nil, err
	}

	//sort.Slice(results.Items, func(i, j int) bool {
	//	return results.Items[i].Index > results.Items[j].Index
	//})

	data := &chunks{
		scores: results.Scores,
	}

	for _, ref := range results.Items {
		for _, chunk := range ref.Chunks {
			data.ids = append(data.ids, chunk.Id)
			data.texts = append(data.texts, chunk.Text)
		}
	}

	return data, nil
}
