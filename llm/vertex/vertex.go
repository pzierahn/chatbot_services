package vertex

import (
	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	"cloud.google.com/go/vertexai/genai"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"os"
)

const (
	localCredentialsFile = "service_account.json"
	projectID            = "brainboost-399710"
	location             = "us-central1"
)

type Client struct {
	ProjectID        string
	Location         string
	predictionClient *aiplatform.PredictionClient
	client           *genai.Client
}

func New(ctx context.Context) (*Client, error) {

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

	client, err := genai.NewClient(
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
		predictionClient: predictionClient,
		client:           client,
	}, nil
}
