package chat

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/utils"
)

type chunks struct {
	ids    []string
	source []string
	texts  []string
	scores map[string]float32
}

func (chunks *chunks) toJSON() string {
	var parts []map[string]string

	for inx := range chunks.source {
		parts = append(parts, map[string]string{
			"SourceID": chunks.source[inx],
			"Text":     chunks.texts[inx],
		})
	}

	byt, _ := json.MarshalIndent(parts, "", "  ")
	return string(byt)
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

	data := &chunks{
		scores: results.Scores,
	}

	for _, ref := range results.Items {
		for _, chunk := range ref.Chunks {
			data.ids = append(data.ids, chunk.Id)
			data.texts = append(data.texts, chunk.Text)

			source := fmt.Sprintf("%s p.%d", utils.GetDocumentTitle(ref.Metadata), chunk.Index+1)
			data.source = append(data.source, source)
		}
	}

	return data, nil
}
