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
	Prompt            string  `json:"prompt"`
	MaxTokensToSample int     `json:"max_tokens_to_sample"`
	Temperature       float32 `json:"temperature"`
}

type ClaudeResponse struct {
	Completion string `json:"completion"`
}

func prettify(obj interface{}) string {
	byt, _ := json.MarshalIndent(obj, "", "  ")
	return string(byt)
}

// main uses the AWS SDK for Go (v2) to create an Amazon Bedrock Runtime client
// and invokes Anthropic Claude 2 inside your account and the chosen region.
// This example uses the default settings specified in your shared credentials
// and config files.
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

	modelId := "anthropic.claude-v2"

	prompt := "Hello, how are you today?"

	// Anthropic Claude requires you to enclose the prompt as follows:
	prefix := "Human: "
	postfix := "\n\nAssistant:"
	wrappedPrompt := prefix + prompt + postfix

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
		ModelId:     aws.String(modelId),
		ContentType: aws.String("application/json"),
		Body:        body,
	})

	if err != nil {
		log.Fatalf("failed to invoke model: %v", err)
	}

	//log.Printf("result.ResultMetadata: %s\n", result.ResultMetadata)

	log.Printf("result.Body: %s\n", result.Body)

	var response ClaudeResponse

	err = json.Unmarshal(result.Body, &response)

	if err != nil {
		log.Fatal("failed to unmarshal", err)
	}
	log.Println("Prompt:", prompt)
	log.Println("Response from Anthropic Claude:", response.Completion)
}

//func listModels() {
//	ctx := context.Background()
//	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
//	if err != nil {
//		log.Println("Couldn't load default configuration. Have you set up your AWS account?")
//		log.Println(err)
//		return
//	}
//
//	bedrockClient := bedrock.NewFromConfig(sdkConfig)
//	result, err := bedrockClient.ListFoundationModels(ctx, &bedrock.ListFoundationModelsInput{})
//	if err != nil {
//		log.Printf("Couldn't list foundation models. Here's why: %v\n", err)
//		return
//	}
//
//	if len(result.ModelSummaries) == 0 {
//		log.Println("There are no foundation models.")
//	}
//
//	for _, modelSummary := range result.ModelSummaries {
//		log.Println(*modelSummary.ModelId)
//	}
//}
