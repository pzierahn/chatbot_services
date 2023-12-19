package vertex

import (
	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	"cloud.google.com/go/vertexai/genai"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"os"
)

type Client struct {
	ProjectID        string
	Location         string
	EmbeddingModel   string
	predictionClient *aiplatform.PredictionClient
	genaiClient      *genai.Client
}

const localCredentialsFile = "brainboost-399710-d789c5991083.json"

func New(ctx context.Context) (*Client, error) {
	projectID := "brainboost-399710"
	location := "us-central1"

	var authOption []option.ClientOption
	if _, err := os.Stat(localCredentialsFile); err == nil {
		localCredentials := option.WithCredentialsFile(localCredentialsFile)
		authOption = append(authOption, localCredentials)
	}

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	predictionClient, err := aiplatform.NewPredictionClient(
		ctx,
		append(authOption, option.WithEndpoint(apiEndpoint))...,
	)
	if err != nil {
		return nil, err
	}

	genaiClient, err := genai.NewClient(
		ctx,
		projectID,
		location,
		authOption...,
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
