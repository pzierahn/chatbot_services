package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/brainboost/database"
	"github.com/sashabaranov/go-openai"
)

type SearchQuery struct {
	UserId     uuid.UUID
	Collection uuid.UUID
	Prompt     string
	Limit      int
	Threshold  float32
}

const (
	embeddingsModel = openai.AdaEmbeddingV2
)

func (server *Server) SearchDocuments(ctx context.Context, query SearchQuery) ([]*database.SearchResult, error) {

	resp, err := server.gpt.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequestStrings{
			Model: embeddingsModel,
			Input: []string{query.Prompt},
			User:  query.UserId.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	_, _ = server.db.CreateUsage(ctx, database.Usage{
		UID:    query.UserId,
		Model:  embeddingsModel.String(),
		Input:  uint32(resp.Usage.PromptTokens),
		Output: uint32(resp.Usage.CompletionTokens),
	})

	return server.db.Search(ctx, database.SearchQuery{
		UserId:     query.UserId,
		Collection: query.Collection,
		Embedding:  resp.Data[0].Embedding,
		Limit:      query.Limit,
		Threshold:  query.Threshold,
	})
}
