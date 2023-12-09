package documents

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"github.com/pzierahn/brainboost/account"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"os"
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

	founding, err := service.account.HasFunding(ctx)
	if err != nil {
		return nil, err
	}

	if !founding {
		return nil, account.NoFundingError()
	}

	out, _ := json.MarshalIndent(query, "", "  ")
	log.Printf("######## query: %s", out)

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

	promptEmbedding := resp.Data[0].Embedding
	_, _ = service.account.CreateUsage(ctx, account.Usage{
		UserId: userID,
		Model:  embeddingsModel.String(),
		Input:  uint32(resp.Usage.PromptTokens),
		Output: uint32(resp.Usage.CompletionTokens),
	})

	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", os.Getenv("PINECONE_KEY"))

	queryResult, err := service.pinecone.Query(ctx, &pinecone_grpc.QueryRequest{
		Queries: []*pinecone_grpc.QueryVector{
			{
				Values: promptEmbedding,
			},
		},
		Filter: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"collectionId": {
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"$eq": {
									Kind: &structpb.Value_StringValue{
										StringValue: query.CollectionId,
									},
								},
							},
						},
					},
				},
			},
		},
		TopK:            200,
		IncludeValues:   false,
		IncludeMetadata: true,
		Namespace:       "documents",
	})
	if err != nil {
		return nil, err
	}

	results := &pb.SearchResults{}

	for _, item := range queryResult.Results[0].Matches {
		if item.Score < query.Threshold {
			break
		}

		if len(results.Items) >= int(query.Limit) {
			break
		}

		doc := &pb.SearchResults_Document{
			Id:         item.Id,
			DocumentId: item.Metadata.Fields["documentId"].GetStringValue(),
			Filename:   item.Metadata.Fields["filename"].GetStringValue(),
			Page:       uint32(item.Metadata.Fields["page"].GetNumberValue()),
			Content:    item.Metadata.Fields["text"].GetStringValue(),
			Score:      item.Score,
		}

		results.Items = append(results.Items, doc)
	}

	return results, nil
}
