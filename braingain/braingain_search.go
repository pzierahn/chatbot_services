package braingain

import (
	"context"
	"github.com/pzierahn/braingain/database"
	"github.com/sashabaranov/go-openai"
	"sort"
)

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

func (chat Chat) Search(ctx context.Context, prompt string) ([]database.ScorePoints, error) {

	embedding, err := chat.createEmbedding(ctx, prompt)
	if err != nil {
		return nil, err
	}

	sources, err := chat.db.SearchEmbedding(ctx, embedding)
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
