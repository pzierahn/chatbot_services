package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"log"
)

type ClaudeMessage struct {
	Role    string    `json:"role,omitempty"`
	Content []Content `json:"content,omitempty"`
}

type ClaudeToolProperty struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

type ClaudeToolInput struct {
	Type       string                        `json:"type,omitempty"`
	Properties map[string]ClaudeToolProperty `json:"properties,omitempty"`
	Required   []string                      `json:"required,omitempty"`
}

type ClaudeTool struct {
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	InputSchema ClaudeToolInput `json:"input_schema,omitempty"`
}

type ClaudeRequest struct {
	AnthropicVersion string          `json:"anthropic_version,omitempty"`
	System           string          `json:"system,omitempty"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	Temperature      float64         `json:"temperature,omitempty"`
	TopP             float64         `json:"top_p,omitempty"`
	TopK             int             `json:"top_k,omitempty"`
	Tools            []ClaudeTool    `json:"tools,omitempty"`
	Messages         []ClaudeMessage `json:"messages,omitempty"`
}

type Content struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`

	// Function Parameters
	ID        string                 `json:"id,omitempty"`
	ToolUseId string                 `json:"tool_use_id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	Input     map[string]interface{} `json:"input,omitempty"`
	Content   string                 `json:"content,omitempty"`
}

type ClaudeUsage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

type ClaudeResponse struct {
	Id         string      `json:"id,omitempty"`
	Model      string      `json:"model,omitempty"`
	Content    []Content   `json:"content,omitempty"`
	Role       string      `json:"role,omitempty"`
	StopReason string      `json:"stop_reason,omitempty"`
	Type       string      `json:"type,omitempty"`
	Usage      ClaudeUsage `json:"usage,omitempty"`
}

const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}

	bedrock := bedrockruntime.NewFromConfig(sdkConfig)

	messages := []ClaudeMessage{{
		Role: ChatMessageRoleUser,
		Content: []Content{{
			Type: "text",
			Text: "Who are Arnold Pitterson and Hugo Alberts von Tahl?",
		}},
	}}

	//
	// First call
	//

	request := ClaudeRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		Messages:         messages,
		System:           "You are a helpful assistant. Quote the sources by \\cite{SourceID}",
		MaxTokens:        256,
		Temperature:      1.0,
		Tools: []ClaudeTool{{
			Name:        "get_sources",
			Description: "Retrieves the sources for the prompt. The prompt should be optimized for embedding retrieval. The tool will return a list of sources in JSON format with the following fields: SourceID, Content.",
			InputSchema: ClaudeToolInput{
				Type: "object",
				Properties: map[string]ClaudeToolProperty{
					"prompt": {
						Type:        "string",
						Description: "The topic for which to retrieve sources. The prompt should be optimized for embedding retrieval.",
					},
				},
				Required: []string{"prompt"},
			},
		}},
	}

	body, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("failed to marshal request, %v", err)
	}

	result, err := bedrock.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String("anthropic.claude-3-5-sonnet-20240620-v1:0"),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		log.Fatalf("failed to invoke model, %v", err)
	}

	var response ClaudeResponse
	err = json.Unmarshal(result.Body, &response)
	if err != nil {
		log.Fatalf("failed to unmarshal response, %v", err)
	}

	byt, _ := json.MarshalIndent(response, "", "  ")
	log.Println("first response:")
	log.Println(string(byt))

	if response.StopReason == "tool_use" {
		messages = append(messages, ClaudeMessage{
			Role:    ChatMessageRoleAssistant,
			Content: response.Content,
		})

		for _, message := range response.Content {
			if message.Type == "tool_use" {
				log.Printf("Tool: %s\n", message.Name)
				log.Printf("Input: %v\n", message.Input)

				sources := []map[string]string{{
					"SourceID": "S1",
					"Content":  "Arnold Pitterson is a fictional character in the book 'The City of Glass' by Paul Auster.",
				}, {
					"SourceID": "S2",
					"Content":  "Hugo Alberts von Tahl was a German philosopher who lived in the 19th century. He is known for his work on the philosophy of language and logic.",
				}}
				sourceByt, _ := json.Marshal(sources)

				messages = append(messages, ClaudeMessage{
					Role: ChatMessageRoleUser,
					Content: []Content{{
						Type:      "tool_result",
						ToolUseId: message.ID,
						Content:   string(sourceByt),
					}},
				})
			}
		}
	}

	byt, _ = json.MarshalIndent(messages, "", "  ")
	log.Println("messages:")
	log.Println(string(byt))

	//
	// Second call
	//

	request = ClaudeRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		Messages:         messages,
		System:           "You are a helpful assistant. Quote the sources by \\cite{SourceID}",
		MaxTokens:        256,
		Tools: []ClaudeTool{{
			Name:        "get_sources",
			Description: "Retrieves the sources for the prompt.",
			InputSchema: ClaudeToolInput{
				Type: "object",
				Properties: map[string]ClaudeToolProperty{
					"prompt": {
						Type:        "string",
						Description: "The prompt for which to retrieve sources.",
					},
				},
				Required: []string{"prompt"},
			},
		}},
	}

	body, err = json.Marshal(request)
	if err != nil {
		log.Fatalf("failed to marshal request, %v", err)
	}

	result, err = bedrock.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String("anthropic.claude-3-5-sonnet-20240620-v1:0"),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		log.Fatalf("failed to invoke model, %v", err)
	}

	err = json.Unmarshal(result.Body, &response)
	if err != nil {
		log.Fatalf("failed to unmarshal response, %v", err)
	}

	byt, _ = json.MarshalIndent(response, "", "  ")
	log.Println("second response:")
	log.Println(string(byt))
}
