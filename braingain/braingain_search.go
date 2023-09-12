package braingain

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/database"
	"github.com/sashabaranov/go-openai"
)

type SearchQuery struct {
	UserId     string
	Collection *uuid.UUID
	Prompt     string
	Limit      int
	Threshold  float32
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

func (chat Chat) Search(ctx context.Context, query SearchQuery) ([]*database.SearchResult, error) {

	embedding, err := chat.createEmbedding(ctx, query.Prompt)
	if err != nil {
		return nil, err
	}

	return chat.db.Search(ctx, database.SearchQuery{
		UserId:     query.UserId,
		Collection: query.Collection,
		Embedding:  embedding,
		Limit:      query.Limit,
		Threshold:  query.Threshold,
	})
}
