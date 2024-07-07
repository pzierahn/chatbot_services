package anthropic

import (
	"encoding/json"
	"github.com/pzierahn/chatbot_services/llm"
)

// transformToClaude converts a list of LLM messages to a list of ClaudeMessages
func transformToClaude(mess []*llm.Message) ([]ClaudeMessage, error) {
	var messages []ClaudeMessage

	for _, msg := range mess {
		var role string
		switch msg.Role {
		case llm.RoleUser:
			role = ChatMessageRoleUser
		case llm.RoleAssistant:
			role = ChatMessageRoleAssistant
		}

		var content []Content
		for _, toolCall := range msg.ToolCalls {
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

		for _, toolResponse := range msg.ToolResponses {
			content = append(content, Content{
				Type:      ContentTypeToolResult,
				ToolUseId: toolResponse.CallID,
				Content:   toolResponse.Content,
			})
		}

		if msg.Content != "" {
			content = append(content, Content{
				Type: ContentTypeText,
				Text: msg.Content,
			})
		}

		messages = append(messages, ClaudeMessage{
			Role:    role,
			Content: content,
		})
	}

	return messages, nil
}
