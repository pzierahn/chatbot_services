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
			err := json.Unmarshal([]byte(call.Function.Arguments), &args)
			if err != nil {
				return nil, err
			}

			history = append(history, &genai.Content{
				Role: RoleModel,
				Parts: []genai.Part{genai.FunctionCall{
					Name: call.Function.Name,
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
						Name:     call.Function.Name,
						Response: response,
					}},
				})
			}
		}
	}

	return history, nil
}
