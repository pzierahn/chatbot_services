package vertex

import (
	"cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/structpb"
)

// GenerateEmbeddings generates embeddings for a prompt
func (client *Client) GenerateEmbeddings(ctx context.Context, prompt string) ([]float32, error) {
	// PredictRequest requires an endpoint, instances, and parameters
	// Endpoint
	base := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models", client.ProjectID, client.Location)
	url := fmt.Sprintf("%s/%s", base, client.EmbeddingModel)

	// Instances: the prompt
	promptValue, err := structpb.NewValue(map[string]interface{}{
		"content": prompt,
	})
	if err != nil {
		return nil, err
	}

	// PredictRequest: create the model prediction request
	req := &aiplatformpb.PredictRequest{
		Endpoint:  url,
		Instances: []*structpb.Value{promptValue},
	}

	// PredictResponse: receive the response from the model
	resp, err := client.predictionClient.Predict(ctx, req)
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

	return embedding, nil
}
