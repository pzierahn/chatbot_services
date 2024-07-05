package openai

import (
	"context"
	"encoding/json"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
	"log"
	"strings"
	"time"
)

func (client *Client) toOpenAIMessage(msg llm.Message) openai.ChatCompletionMessage {
	var role string
	switch msg.Role {
	case llm.MessageTypeUser:
		role = openai.ChatMessageRoleUser
	case llm.MessageTypeAssistant:
		role = openai.ChatMessageRoleAssistant
	case llm.MessageTypeTool:
		role = openai.ChatMessageRoleTool
	}

	message := openai.ChatCompletionMessage{
		Role:      role,
		Content:   msg.Content,
		ToolCalls: make([]openai.ToolCall, len(msg.ToolCalls)),
	}

	for inx, call := range msg.ToolCalls {
		message.ToolCalls[inx] = openai.ToolCall{
			ID: call.Id,
			Function: openai.FunctionCall{
				Name:      call.Function.Name,
				Arguments: call.Function.Arguments,
			},
		}
	}

	return message
}

func (client *Client) Completion(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	var messages []openai.ChatCompletionMessage

	if req.SystemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.SystemPrompt,
		})
	}

	for _, msg := range req.Messages {
		messages = append(messages, client.toOpenAIMessage(*msg))
	}

	model, _ := strings.CutPrefix(req.Model, modelPrefix)

	request := openai.ChatCompletionRequest{
		Model:       model,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		MaxTokens:   req.MaxTokens,
		Messages:    messages,
		Tools:       client.getTools(),
		N:           1,
		User:        req.UserId,
	}

	byt, _ := json.MarshalIndent(request, "", "  ")
	log.Println("request:", string(byt))

	resp, err := client.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, err
	}

	//
	// The model wants to call a tools
	//

	if resp.Choices[0].FinishReason == openai.FinishReasonToolCalls {
		messages = append(messages, resp.Choices[0].Message)

		for _, tool := range resp.Choices[0].Message.ToolCalls {
			function := tool.Function

			response, err := client.callTool(ctx, function.Name, function.Arguments)
			if err != nil {
				return nil, err
			}

			message := openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    response,
				ToolCallID: tool.ID,
			}

			messages = append(messages, message)
		}

		request = openai.ChatCompletionRequest{
			Model:       model,
			Temperature: req.Temperature,
			MaxTokens:   req.MaxTokens,
			Messages:    messages,
			N:           1,
			Tools:       client.getTools(),
			User:        req.UserId,
		}

		resp, err = client.client.CreateChatCompletion(ctx, request)
		if err != nil {
			return nil, err
		}
	}

	return &llm.CompletionResponse{
		Message: &llm.Message{
			Role:      llm.MessageTypeAssistant,
			Content:   resp.Choices[0].Message.Content,
			Timestamp: time.Now(),
		},
		Usage: llm.ModelUsage{
			UserId:       req.UserId,
			Model:        resp.Model,
			InputTokens:  resp.Usage.PromptTokens,
			OutputTokens: resp.Usage.CompletionTokens,
		},
	}, nil
}
