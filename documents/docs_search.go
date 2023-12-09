package documents

import (
	"context"
	"github.com/google/uuid"
	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"github.com/pzierahn/brainboost/account"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
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
				"userId": {
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"$eq": {
									Kind: &structpb.Value_StringValue{
										StringValue: userId,
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
