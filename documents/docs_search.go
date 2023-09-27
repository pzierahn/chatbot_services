package documents

import (
	"context"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"github.com/pzierahn/brainboost/account"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
)

type SearchQuery struct {
	UserId     uuid.UUID
	Collection uuid.UUID
	Prompt     string
	Limit      int
	Threshold  float32
}

func (service *Service) Search(ctx context.Context, query *pb.SearchQuery) (*pb.SearchResults, error) {

	userID, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := service.gpt.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequestStrings{
			Model: embeddingsModel,
			Input: []string{query.Query},
			User:  userID.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	embedding := resp.Data[0].Embedding
	_, _ = service.account.CreateUsage(ctx, account.Usage{
		UserID: userID,
		Model:  embeddingsModel.String(),
		Input:  uint32(resp.Usage.PromptTokens),
		Output: uint32(resp.Usage.CompletionTokens),
	})

	rows, err := service.db.Query(
		ctx,
		`select em.id, document_id, filename, page, text, (1 - (embedding <=> $1)) AS score
			from document_embeddings as em join documents as doc on doc.id = em.document_id
			where (1 - (embedding <=> $1)) >= $2 and doc.user_id = $3 and doc.collection_id = $4
			ORDER BY score DESC
		 	LIMIT $5`,
		pgvector.NewVector(embedding),
		query.Threshold,
		userID,
		query.CollectionID,
		query.Limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := &pb.SearchResults{}
	for rows.Next() {
		var doc pb.SearchResults_Document

		err = rows.Scan(
			&doc.Id,
			&doc.DocumentID,
			&doc.Filename,
			&doc.Page,
			&doc.Content,
			&doc.Score)
		if err != nil {
			return nil, err
		}

		results.Items = append(results.Items, &doc)
	}

	return results, nil
}
