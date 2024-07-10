package openai

import (
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
)

// messagesToOpenAI converts messages from the LLM format to the OpenAI format.
func messagesToOpenAI(input []*llm.Message) []openai.ChatCompletionMessage {
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
						Name:      call.Name,
						Arguments: call.Arguments,
					},
				}
			}

			messages = append(messages, message)
		}
	}

	return messages
}

// openaiToMessages converts messages from the OpenAI format to the LLM format.
func openaiToMessages(input []openai.ChatCompletionMessage) []*llm.Message {
	var messages []*llm.Message

	for _, msg := range input {
		switch msg.Role {
		case openai.ChatMessageRoleUser:
			messages = append(messages, &llm.Message{
				Role:    llm.RoleUser,
				Content: msg.Content,
			})
		case openai.ChatMessageRoleAssistant:
			message := &llm.Message{
				Role:    llm.RoleAssistant,
				Content: msg.Content,
			}

			for _, call := range msg.ToolCalls {
				message.ToolCalls = append(message.ToolCalls, llm.ToolCall{
					CallID:    call.ID,
					Name:      call.Function.Name,
					Arguments: call.Function.Arguments,
				})
			}

			messages = append(messages, message)
		case openai.ChatMessageRoleTool:
			// Check if the last message is a user message and create a new one if not
			if len(messages) == 0 || messages[len(messages)-1].Role != llm.RoleUser {
				messages = append(messages, &llm.Message{
					Role: llm.RoleUser,
				})
			}

			idx := len(messages) - 1
			messages[idx].ToolResponses = append(messages[idx].ToolResponses, llm.ToolResponse{
				CallID:  msg.ToolCallID,
				Content: msg.Content,
			})
		}
	}

	return messages
}
