package documents

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/llm"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/pzierahn/brainboost/vectordb"
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

	resp, err := service.embeddings.CreateEmbeddings(ctx, &llm.EmbeddingRequest{
		Input:  query.Query,
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	_, _ = service.account.CreateUsage(ctx, account.Usage{
		UserId: userId,
		Model:  embeddingsModel.String(),
		Input:  uint32(resp.Tokens),
	})

	results, err := service.vectorDB.Search(vectordb.SearchQuery{
		UserId:       userId,
		CollectionId: query.CollectionId,
		Vector:       resp.Data,
		Limit:        int(query.Limit),
		Threshold:    query.Threshold,
	})

	return results, err
}
