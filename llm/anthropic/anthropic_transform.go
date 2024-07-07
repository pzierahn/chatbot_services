package anthropic

import (
	"encoding/json"
	"github.com/pzierahn/chatbot_services/llm"
)

// transformToClaude converts a list of LLM messages to a list of ClaudeMessages
func transformToClaude(messages []*llm.Message) ([]ClaudeMessage, error) {
	var claudeMessages []ClaudeMessage

	for _, message := range messages {
		var role string
		switch message.Role {
		case llm.RoleUser:
			role = ChatMessageRoleUser
		case llm.RoleAssistant:
			role = ChatMessageRoleAssistant
		}

		var content []Content
		for _, toolCall := range message.ToolCalls {
			var args map[string]interface{}
			err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
			if err != nil {
				return nil, err
			}

			content = append(content, Content{
				Type:  ContentTypeToolUse,
				ID:    toolCall.CallID,
				Name:  toolCall.Function.Name,
				Input: args,
			})
		}

		for _, toolResponse := range message.ToolResponses {
			content = append(content, Content{
				Type:      ContentTypeToolResult,
				ToolUseId: toolResponse.CallID,
				Content:   toolResponse.Content,
			})
		}

		if message.Content != "" {
			content = append(content, Content{
				Type: ContentTypeText,
				Text: message.Content,
			})
		}

		claudeMessages = append(claudeMessages, ClaudeMessage{
			Role:    role,
			Content: content,
		})
	}

	return claudeMessages, nil
}

// claudeToMessages converts a list of ClaudeMessages to a list of LLM messages
func claudeToMessages(messages []ClaudeMessage) ([]*llm.Message, error) {
	var llmMessages []*llm.Message

	for _, message := range messages {
		var role string
		switch message.Role {
		case ChatMessageRoleUser:
			role = llm.RoleUser
		case ChatMessageRoleAssistant:
			role = llm.RoleAssistant
		}

		llmMessage := &llm.Message{
			Role: role,
		}

		for _, content := range message.Content {
			switch content.Type {
			case ContentTypeToolUse:
				args, err := json.Marshal(content.Input)
				if err != nil {
					return nil, err
				}

				llmMessage.ToolCalls = append(llmMessage.ToolCalls, llm.ToolCall{
					CallID: content.ID,
					Function: llm.Function{
						Name:      content.Name,
						Arguments: string(args),
					},
				})
			case ContentTypeToolResult:
				llmMessage.ToolResponses = append(llmMessage.ToolResponses, llm.ToolResponse{
					CallID:  content.ToolUseId,
					Content: content.Content,
				})
			case ContentTypeText:
				llmMessage.Content = content.Text
			}
		}

		llmMessages = append(llmMessages, llmMessage)
	}

	return llmMessages, nil
}
