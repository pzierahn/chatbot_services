package braingain

import (
	"context"
	"github.com/pzierahn/braingain/database"
	"github.com/sashabaranov/go-openai"
	"sort"
)

type SearchQuery struct {
	Prompt    string
	Limit     int
	Threshold float32
}

func (chat Chat) createEmbedding(ctx context.Context, prompt string) ([]float32, error) {
	resp, err := chat.gpt.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequestStrings{
			Model: openai.AdaEmbeddingV2,
			Input: []string{prompt},
		},
	)

	if err != nil {
		return nil, err
	}

	return resp.Data[0].Embedding, nil
}

func (chat Chat) Search(ctx context.Context, query SearchQuery) ([]database.ScorePoints, error) {

	embedding, err := chat.createEmbedding(ctx, query.Prompt)
	if err != nil {
		return nil, err
	}

	sources, err := chat.db.SearchEmbedding(ctx, database.SearchQuery{
		Embedding: embedding,
		Limit:     query.Limit,
		Threshold: query.Threshold,
	})
	if err != nil {
		return nil, err
	}

	sort.SliceStable(sources, func(i, j int) bool {
		return sources[i].Page < sources[j].Page
	})
	sort.SliceStable(sources, func(i, j int) bool {
		return sources[i].Source.String() < sources[j].Source.String()
	})

	return sources, nil
}
