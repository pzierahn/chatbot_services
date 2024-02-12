package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
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

	log.Printf("Using AWS region: %s\n", region)

	sdkConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Println("Couldn't load default configuration. Have you set up your AWS account?")
		log.Println(err)
		return
	}

	client := bedrockruntime.NewFromConfig(sdkConfig)

	prompt := "Hello, how are you today?"
	wrappedPrompt := "Human: " + prompt + "\n\nAssistant:"
	request := ClaudeRequest{
		Prompt:            wrappedPrompt,
		MaxTokensToSample: 200,
		Temperature:       0.5,
	}

	body, err := json.Marshal(request)
	if err != nil {
		log.Fatal("failed to marshal request", err)
	}

	result, err := client.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String("anthropic.claude-v2"),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		log.Fatalf("failed to invoke model: %v", err)
	}

	// The metadata are not accessible in the results.
	// The map keys are always {} and there is no way to iterate over the underlying values map.
	log.Printf("result.ResultMetadata: %v\n", result.ResultMetadata)

	log.Printf("result.Body: %s\n", result.Body)

	var response ClaudeResponse

	err = json.Unmarshal(result.Body, &response)

	if err != nil {
		log.Fatal("failed to unmarshal", err)
	}
	log.Println("Response from Anthropic Claude:", response.Completion)
}
