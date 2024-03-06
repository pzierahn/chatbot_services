package vertex

import (
	"cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
)

// CreateEmbeddings generates embeddings for a text
func (client *Client) CreateEmbeddings(ctx context.Context, req *llm.EmbeddingRequest) (*llm.EmbeddingResponse, error) {
	base := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models", client.ProjectID, client.Location)
	url := fmt.Sprintf("%s/%s", base, client.EmbeddingModel)

	promptValue, err := structpb.NewValue(map[string]interface{}{
		//"text": req.Input,
		"image": map[string]interface{}{
			"bytesBase64Encoded": req.Input,
		},
	})
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	resp, err := client.predictionClient.Predict(ctx, &aiplatformpb.PredictRequest{
		Endpoint:  url,
		Instances: []*structpb.Value{promptValue},
	})
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	if len(resp.Predictions) == 0 {
		return nil, fmt.Errorf("no predictions")
	}

	embeddings := resp.Predictions[0].GetStructValue().AsMap()

	var embedding []float32
	textEmbedding, ok := embeddings["textEmbedding"].([]interface{})
	if ok {
		embedding = make([]float32, len(textEmbedding))
		for inx, value := range textEmbedding {
			embedding[inx] = float32(value.(float64))
		}
	}

	imageEmbedding, ok := embeddings["imageEmbedding"].([]interface{})
	if ok {
		embedding = make([]float32, len(imageEmbedding))
		for inx, value := range imageEmbedding {
			embedding[inx] = float32(value.(float64))
		}
	}

	if len(embedding) == 0 {
		return nil, fmt.Errorf("values not found")
	}

	// TODO: Count tokens

	return &llm.EmbeddingResponse{
		Data: embedding,
	}, nil
}

// GetEmbeddingModelName GetModelName returns the name of the model
func (client *Client) GetEmbeddingModelName() string {
	return client.EmbeddingModel
}
