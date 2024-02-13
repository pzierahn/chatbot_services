package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/transport/http"
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

	deserializeMiddleware := middleware.DeserializeMiddlewareFunc(
		"myDeserializeMiddleware",
		func(ctx context.Context, input middleware.DeserializeInput, next middleware.DeserializeHandler) (middleware.DeserializeOutput, middleware.Metadata, error) {
			output, metadata, err := next.HandleDeserialize(ctx, input)
			log.Printf("output.RawResponse: %v", output.RawResponse)

			if resp, ok := output.RawResponse.(*http.Response); ok {
				log.Printf("resp.Header: %v", resp.Header)
			}
			return output, metadata, err
		})

	client := bedrockruntime.NewFromConfig(sdkConfig, func(options *bedrockruntime.Options) {
		options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
			return stack.Deserialize.Insert(deserializeMiddleware, "OperationDeserializer", middleware.After)
		})
	})

	request := ClaudeRequest{
		Prompt: "\n\nHuman: I have a little green rectangular object in a yellow box\n\n" +
			"Assistant: ",
		MaxTokensToSample: 100,
	}

	body, err := json.Marshal(request)
	if err != nil {
		log.Fatal("failed to marshal request", err)
	}

	result, err := client.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String("anthropic.claude-v2"),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
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
