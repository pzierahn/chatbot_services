package documents

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/search"
	pb "github.com/pzierahn/chatbot_services/services/proto"
)

type SearchQuery struct {
	UserId     uuid.UUID
	Collection uuid.UUID
	Prompt     string
	Limit      int
	Threshold  float32
}

func (service *Service) Search(ctx context.Context, query *pb.SearchQuery) (*pb.SearchResults, error) {

	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	searchResults, err := service.SearchIndex.Search(ctx, search.Query{
		UserId:       userId,
		CollectionId: query.CollectionId,
		Query:        query.Text,
		Limit:        query.Limit,
		Threshold:    query.Threshold,
	})
	if err != nil {
		return nil, err
	}

	results := &pb.SearchResults{
		DocumentNames: make(map[string]string),
		Scores:        make(map[string]float32),
	}

	for _, vector := range searchResults.Results {
		results.Chunks = append(results.Chunks, &pb.Chunk{
			Id:      vector.Id,
			Text:    vector.Text,
			Postion: vector.Position,
		})

		results.DocumentNames[vector.DocumentId] = ""
		results.Scores[vector.Id] = vector.Score
	}

	docIds := make([]uuid.UUID, 0)
	for docId := range results.DocumentNames {
		docIds = append(docIds, uuid.MustParse(docId))
	}

	docs, err := service.Database.GetDocumentMeta(ctx, userId, docIds...)
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		results.DocumentNames[doc.Id.String()] = doc.Name
	}

	return results, nil
}
