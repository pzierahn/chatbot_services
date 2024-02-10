package documents

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/vectordb"
)

type SearchQuery struct {
	UserId     uuid.UUID
	Collection uuid.UUID
	Prompt     string
	Limit      int
	Threshold  float32
}

func (service *Service) Search(ctx context.Context, query *pb.SearchQuery) (*pb.SearchResults, error) {

	userId, err := service.auth.Verify(ctx)
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

	vectors, err := service.vectorDB.Search(vectordb.SearchQuery{
		UserId:       userId,
		CollectionId: query.CollectionId,
		Vector:       resp.Data,
		Limit:        int(query.Limit),
		Threshold:    query.Threshold,
	})
	if err != nil {
		return nil, err
	}

	results := &pb.SearchResults{
		Items: make([]*pb.SearchResults_Document, len(vectors)),
	}
	for inx, vector := range vectors {
		results.Items[inx] = &pb.SearchResults_Document{
			Id:         vector.Id,
			DocumentId: vector.DocumentId,
			Filename:   vector.Filename,
			Content:    vector.Text,
			Page:       vector.Page,
			Score:      vector.Score,
		}
	}

	return results, nil
}
