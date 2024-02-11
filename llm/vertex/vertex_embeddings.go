package vertex

import (
	"cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	"google.golang.org/protobuf/types/known/structpb"
)

// CreateEmbeddings generates embeddings for a text
func (client *Client) CreateEmbeddings(ctx context.Context, req *llm.EmbeddingRequest) (*llm.EmbeddingResponse, error) {
	base := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models", client.ProjectID, client.Location)
	url := fmt.Sprintf("%s/%s", base, client.EmbeddingModel)

	promptValue, err := structpb.NewValue(map[string]interface{}{
		"content": req.Input,
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.predictionClient.Predict(ctx, &aiplatformpb.PredictRequest{
		Endpoint:  url,
		Instances: []*structpb.Value{promptValue},
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Predictions) == 0 {
		return nil, fmt.Errorf("no predictions")
	}

	pred := resp.Predictions[0].GetStructValue().AsMap()
	embeddings, ok := pred["embeddings"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("embeddings not found")
	}

	values, ok := embeddings["values"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("values not found")
	}

	embedding := make([]float32, len(values))
	for i, v := range values {
		embedding[i] = float32(v.(float64))
	}

	meta := resp.Metadata.GetStructValue().AsMap()
	charCount, ok := meta["billableCharacterCount"].(float64)
	if !ok {
		return nil, fmt.Errorf("billableCharacterCount not found")
	}

	// TODO: Add usage tracking

	return &llm.EmbeddingResponse{
		Data:   embedding,
		Tokens: int(charCount),
	}, nil
}

// GetModelName returns the name of the model
func (client *Client) GetModelName() string {
	return client.EmbeddingModel
}
