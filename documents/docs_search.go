package documents

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/brainboost/account"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/pzierahn/brainboost/vectordb"
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

	userId, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	funding, err := service.account.HasFunding(ctx)
	if err != nil {
		return nil, err
	}

	if !funding {
		return nil, account.NoFundingError()
	}

	resp, err := service.gpt.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequestStrings{
			Model: embeddingsModel,
			Input: []string{query.Query},
			User:  userId,
		},
	)
	if err != nil {
		return nil, err
	}

	promptEmbedding := resp.Data[0].Embedding
	_, _ = service.account.CreateUsage(ctx, account.Usage{
		UserId: userId,
		Model:  embeddingsModel.String(),
		Input:  uint32(resp.Usage.PromptTokens),
		Output: uint32(resp.Usage.CompletionTokens),
	})

	results, err := service.vectorDB.Search(vectordb.SearchQuery{
		UserId:       userId,
		CollectionId: query.CollectionId,
		Vector:       promptEmbedding,
		Limit:        int(query.Limit),
		Threshold:    query.Threshold,
	})

	return results, err
}
