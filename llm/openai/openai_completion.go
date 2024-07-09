package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
	"strings"
)

func (client *Client) Completion(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	var messages []openai.ChatCompletionMessage

	if req.SystemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.SystemPrompt,
		})
	}

	messages = append(messages, messagesToOpenAI(req.Messages)...)
	model, _ := strings.CutPrefix(req.Model, modelPrefix)

	tools := toolConverter(req.Tools)

	request := openai.ChatCompletionRequest{
		Model:       model,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		MaxTokens:   req.MaxTokens,
		Messages:    messages,
		Tools:       tools.toOpenAI(),
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
		InputTokens:  uint32(resp.Usage.PromptTokens),
		OutputTokens: uint32(resp.Usage.CompletionTokens),
	}

	if resp.Choices[0].FinishReason == openai.FinishReasonToolCalls {
		//
		// The model wants to call tools
		//

		request.Messages = append(request.Messages, resp.Choices[0].Message)

		for _, tool := range resp.Choices[0].Message.ToolCalls {
			name := tool.Function.Name
			arguments := tool.Function.Arguments

			function, ok := tools.getFunction(name)
			if !ok {
				return nil, fmt.Errorf("unknown tool function: %s", name)
			}

			var input map[string]interface{}
			if arguments != "" {
				err = json.Unmarshal([]byte(arguments), &input)
				if err != nil {
					return nil, fmt.Errorf("invalid tool arguments: %s", arguments)
				}
			}

			response, err := function(ctx, input)
			if err != nil {
				return nil, err
			}

			message := openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    response,
				ToolCallID: tool.ID,
			}

			request.Messages = append(request.Messages, message)
		}

		resp, err = client.client.CreateChatCompletion(ctx, request)
		if err != nil {
			return nil, err
		}

		// Add the tool usage to the model usage
		usage.InputTokens += uint32(resp.Usage.PromptTokens)
		usage.OutputTokens += uint32(resp.Usage.CompletionTokens)
	}

	thread := openaiToMessages(request.Messages)
	thread = append(thread, &llm.Message{
		Role:    llm.RoleAssistant,
		Content: strings.TrimSpace(resp.Choices[0].Message.Content),
	})

	return &llm.CompletionResponse{
		Messages: thread,
		Usage:    usage,
	}, nil
}
