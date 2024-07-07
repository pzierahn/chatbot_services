package openai

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
	"strings"
)

func (client *Client) toOpenAIMessages(input []*llm.Message) []openai.ChatCompletionMessage {
	var messages []openai.ChatCompletionMessage

	for _, msg := range input {
		switch msg.Role {
		case llm.RoleUser:

			if msg.Content != "" {
				messages = append(messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: msg.Content,
				})
			}

			for _, response := range msg.ToolResponses {
				messages = append(messages, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    response.Content,
					ToolCallID: response.CallID,
				})
			}
		case llm.RoleAssistant:
			message := openai.ChatCompletionMessage{
				Role:      openai.ChatMessageRoleAssistant,
				Content:   msg.Content,
				ToolCalls: make([]openai.ToolCall, len(msg.ToolCalls)),
			}

			for inx, call := range msg.ToolCalls {
				message.ToolCalls[inx] = openai.ToolCall{
					ID:   call.CallID,
					Type: openai.ToolTypeFunction,
					Function: openai.FunctionCall{
						Name:      call.Function.Name,
						Arguments: call.Function.Arguments,
					},
				}
			}

			messages = append(messages, message)
		}
	}

	return messages
}

func (client *Client) Completion(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	var messages []openai.ChatCompletionMessage

	if req.SystemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.SystemPrompt,
		})
	}

	messages = append(messages, client.toOpenAIMessages(req.Messages)...)

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

	resp, err := client.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, err
	}

	usage := llm.ModelUsage{
		UserId:       req.UserId,
		Model:        resp.Model,
		InputTokens:  resp.Usage.PromptTokens,
		OutputTokens: resp.Usage.CompletionTokens,
	}

	if resp.Choices[0].FinishReason == openai.FinishReasonToolCalls {
		//
		// The model wants to call tools
		//

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

		// Add the tool usage to the model usage
		usage.OutputTokens += resp.Usage.CompletionTokens
		usage.InputTokens += resp.Usage.PromptTokens
	}

	return &llm.CompletionResponse{
		Message: &llm.Message{
			Role:    llm.RoleAssistant,
			Content: resp.Choices[0].Message.Content,
		},
		Usage: usage,
	}, nil
}
