package server

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

const (
	embeddingsModel = openai.AdaEmbeddingV2
)

func (server *Server) createEmbedding(ctx context.Context, prompt string) ([]float32, error) {
	resp, err := server.gpt.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequestStrings{
			Model: embeddingsModel,
			Input: []string{prompt},
		},
	)

	if err != nil {
		return nil, err
	}

	return resp.Data[0].Embedding, nil
}

func (server *Server) SearchDocuments(ctx context.Context, query SearchQuery) ([]*database.SearchResult, error) {

	embedding, err := server.createEmbedding(ctx, query.Prompt)
	if err != nil {
		return nil, err
	}

	return server.db.Search(ctx, database.SearchQuery{
		UserId:     query.UserId,
		Collection: query.Collection,
		Embedding:  embedding,
		Limit:      query.Limit,
		Threshold:  query.Threshold,
	})
}
