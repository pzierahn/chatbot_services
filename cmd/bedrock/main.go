package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/pzierahn/chatbot_services/utils"
	"log"
)

const region = "us-east-1"

type ClaudeRequest struct {
	Prompt            string   `json:"prompt"`
	MaxTokensToSample int      `json:"max_tokens_to_sample"`
	Temperature       float64  `json:"temperature,omitempty"`
	TopP              float64  `json:"top_p,omitempty"`
	TopK              int      `json:"top_k,omitempty"`
	StopSequences     []string `json:"stop_sequences,omitempty"`
}

type ClaudeResponse struct {
	Completion string `json:"completion"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	sdkConfig, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithClientLogMode(aws.LogResponseWithBody))
	if err != nil {
		log.Printf("Couldn't load default configuration: %v\n", err)
		return
	}

	client := bedrockruntime.NewFromConfig(sdkConfig)

	request := ClaudeRequest{
		Prompt: "\n\nHuman: I have a little green rectangular object in a yellow box\n\n" +
			"Assistant: ",
		MaxTokensToSample: 100,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		log.Fatal("failed to marshal request", err)
	}

	result, err := client.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String("anthropic.claude-v2"),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        jsonBody,
	})
	if err != nil {
		log.Fatalf("failed to invoke model: %v", err)
	}

	var response ClaudeResponse
	err = json.Unmarshal(result.Body, &response)
	if err != nil {
		log.Fatal("failed to unmarshal", err)
	}
	log.Println("Response from bedrock:", utils.Prettify(response))
}
