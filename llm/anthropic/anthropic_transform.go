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
