package anthropic

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type Client struct {
	bedrock *bedrockruntime.Client
}

const region = "us-east-1"

func New() (*Client, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return &Client{
		bedrock: bedrockruntime.NewFromConfig(sdkConfig),
	}, nil
}
