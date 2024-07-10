package vertex

import (
	"cloud.google.com/go/vertexai/genai"
	"encoding/json"
	"github.com/pzierahn/chatbot_services/llm"
)

// transformToHistory transforms a list of messages to a list of history items
func transformToHistory(messages []*llm.Message) ([]*genai.Content, error) {
	var history []*genai.Content

	for inx, msg := range messages {
		var role string

		if msg.Role == llm.RoleUser {
			role = RoleUser
		} else {
			role = RoleModel
		}

		if msg.Content != "" {
			history = append(history, &genai.Content{
				Role:  role,
				Parts: []genai.Part{genai.Text(msg.Content)},
			})
		}

		for iny, call := range msg.ToolCalls {
			var args map[string]interface{}
			err := json.Unmarshal([]byte(call.Arguments), &args)
			if err != nil {
				return nil, err
			}

			history = append(history, &genai.Content{
				Role: RoleModel,
				Parts: []genai.Part{genai.FunctionCall{
					Name: call.Name,
					Args: args,
				}},
			})

			// Check if the next message is a tool response
			if inx+1 < len(messages) && len(messages[inx+1].ToolResponses) > iny {
				toolResponse := messages[inx+1].ToolResponses[iny]

				// Parse the tool response
				var response map[string]interface{}
				err = json.Unmarshal([]byte(toolResponse.Content), &response)
				if err != nil {
					return nil, err
				}

				// Add the tool response to the history
				history = append(history, &genai.Content{
					Role: RoleUser,
					Parts: []genai.Part{genai.FunctionResponse{
						Name:     call.Name,
						Response: response,
					}},
				})
			}
		}
	}

	return history, nil
}

// transformToMessages transforms a list of history items to a list of messages
func transformToMessages(history []*genai.Content) ([]*llm.Message, error) {
	var messages []*llm.Message

	for _, content := range history {
		var role string

		if content.Role == RoleUser {
			role = llm.RoleUser
		} else {
			role = llm.RoleAssistant
		}

		message := &llm.Message{
			Role: role,
		}

		for _, part := range content.Parts {
			if txt, ok := part.(genai.Text); ok {
				message.Content = string(txt)
			} else if call, ok := part.(genai.FunctionCall); ok {
				args, err := json.Marshal(call.Args)
				if err != nil {
					return nil, err
				}

				message.ToolCalls = append(message.ToolCalls, llm.ToolCall{
					Name:      call.Name,
					Arguments: string(args),
				})
			} else if response, ok := part.(genai.FunctionResponse); ok {
				resp, err := json.Marshal(response.Response)
				if err != nil {
					return nil, err
				}

				message.ToolResponses = append(message.ToolResponses, llm.ToolResponse{
					Content: string(resp),
				})
			}
		}

		messages = append(messages, message)
	}

	return messages, nil
}
