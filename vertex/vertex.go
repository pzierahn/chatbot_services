package vertex

import (
	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	"cloud.google.com/go/vertexai/genai"
	"context"
	"fmt"
	"google.golang.org/api/option"
)

type Client struct {
	ProjectID        string
	Location         string
	EmbeddingModel   string
	GenerationModel  string
	predictionClient *aiplatform.PredictionClient
	genaiClient      *genai.Client
}

func New(ctx context.Context) (*Client, error) {
	projectID := "brainboost-399710"
	location := "us-central1"

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)

	predictionClient, err := aiplatform.NewPredictionClient(
		ctx,
		option.WithEndpoint(apiEndpoint),
		option.WithCredentialsFile("brainboost-399710-d789c5991083.json"),
	)
	if err != nil {
		return nil, err
	}

	genaiClient, err := genai.NewClient(
		ctx,
		projectID,
		location,
		option.WithCredentialsFile("brainboost-399710-d789c5991083.json"),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		ProjectID:        projectID,
		Location:         location,
		EmbeddingModel:   "textembedding-gecko-multilingual",
		predictionClient: predictionClient,
		genaiClient:      genaiClient,
	}, nil
}
