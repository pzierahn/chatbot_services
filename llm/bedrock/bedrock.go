package bedrock

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/pzierahn/chatbot_services/llm"
)

type Client struct {
	bedrock *bedrockruntime.Client
	usage   llm.Usage
}

const region = "us-east-1"

func New(usage llm.Usage) (*Client, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return &Client{
		bedrock: bedrockruntime.NewFromConfig(sdkConfig),
		usage:   usage,
	}, nil
}
