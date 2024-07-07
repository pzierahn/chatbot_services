package anthropic

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

func (client *Client) invokeRequest(model string, req *ClaudeRequest) (*ClaudeResponse, error) {
	body, _ := json.Marshal(req)
	result, err := client.bedrock.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(model),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return nil, err
	}

	var response ClaudeResponse
	err = json.Unmarshal(result.Body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *Client) Completion(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {

	messages, err := transformToClaude(req.Messages)
	if err != nil {
		return nil, err
	}

	request := ClaudeRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		Messages:         messages,
		System:           req.SystemPrompt,
		MaxTokens:        req.MaxTokens,
		TopP:             req.TopP,
		Temperature:      req.Temperature,
		Tools:            client.getTools(),
	}

	response, err := client.invokeRequest(req.Model, &request)
	if err != nil {
		return nil, err
	}

	usage := llm.ModelUsage{
		UserId:       req.UserId,
		Model:        response.Model,
		InputTokens:  response.Usage.InputTokens,
		OutputTokens: response.Usage.OutputTokens,
	}

	loops := 0
	for response.StopReason == ContentTypeToolUse && loops < 6 {
		messages = append(messages, ClaudeMessage{
			Role:    ChatMessageRoleAssistant,
			Content: response.Content,
		})

		for _, message := range response.Content {
			if message.Type == ContentTypeToolUse {
				result, err := client.callTool(ctx, message.Name, message.Input)
				if err != nil {
					return nil, err
				}

				messages = append(messages, ClaudeMessage{
					Role: ChatMessageRoleUser,
					Content: []Content{{
						Type:      ContentTypeToolResult,
						ToolUseId: message.ID,
						Content:   result,
					}},
				})
			}
		}

		request.Messages = messages
		response, err = client.invokeRequest(req.Model, &request)
		if err != nil {
			return nil, err
		}

		usage.OutputTokens += response.Usage.OutputTokens
		usage.InputTokens += response.Usage.InputTokens

		loops++
	}

	return &llm.CompletionResponse{
		Message: &llm.Message{
			Role:    llm.RoleAssistant,
			Content: strings.TrimSpace(response.Content[0].Text),
		},
		Usage: usage,
	}, nil
}
