package vertex

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"encoding/json"
	"errors"
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

const (
	RoleUser  = "user"
	RoleModel = "model"
)

func (client *Client) transformToHistory(messages []*llm.Message) ([]*genai.Content, error) {
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
			if inx+1 < len(messages) && messages[inx+1].Role == llm.RoleTool {
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

func (client *Client) Completion(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if len(req.Messages) == 0 {
		return nil, errors.New("no messages")
	}

	modelName, _ := strings.CutPrefix(req.Model, modelPrefix)

	outputTokens := int32(req.MaxTokens)

	model := client.client.GenerativeModel(modelName)
	model.TopP = &req.TopP
	model.Temperature = &req.Temperature
	model.MaxOutputTokens = &outputTokens
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(req.SystemPrompt)},
	}
	model.Tools = client.getTools()

	chat := model.StartChat()

	// Transform the messages to a history
	history, err := client.transformToHistory(req.Messages)
	if err != nil {
		return nil, err
	}

	// Remove the last message from the history, because the last message needs to be sent to the model
	chat.History = history[:len(history)-1]
	gen, err := chat.SendMessage(ctx, history[len(history)-1].Parts...)
	if err != nil {
		return nil, err
	}

	// Check if the model returned a response
	if len(gen.Candidates) == 0 || len(gen.Candidates[0].Content.Parts) == 0 {
		return nil, nil
	}

	usage := llm.ModelUsage{
		UserId: req.UserId,
		Model:  modelName,
	}

	if gen.UsageMetadata != nil {
		usage.InputTokens = int(gen.UsageMetadata.PromptTokenCount)
		usage.OutputTokens = int(gen.UsageMetadata.CandidatesTokenCount)
	}

	if fun, ok := gen.Candidates[0].Content.Parts[0].(genai.FunctionCall); ok {
		// Add the function call to the history
		history = append(history, &genai.Content{
			Role:  RoleModel,
			Parts: []genai.Part{fun},
		})

		// Call the function to get the result
		results, err := client.callTool(ctx, fun.Name, fun.Args)
		if err != nil {
			return nil, err
		}

		functionResults := genai.FunctionResponse{
			Name: fun.Name,
			Response: map[string]any{
				"content": results,
			},
		}
		history = append(history, &genai.Content{
			Role:  RoleUser,
			Parts: []genai.Part{functionResults},
		})
		chat.History = history[:len(history)-1]

		gen, err = chat.SendMessage(ctx, functionResults)
		if err != nil {
			return nil, err
		}

		if gen.UsageMetadata != nil {
			usage.InputTokens += int(gen.UsageMetadata.PromptTokenCount)
			usage.OutputTokens += int(gen.UsageMetadata.CandidatesTokenCount)
		}
	}

	txt, ok := gen.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, nil
	}

	return &llm.CompletionResponse{
		Message: &llm.Message{
			Role:    llm.RoleAssistant,
			Content: string(txt),
		},
		Usage: usage,
	}, nil
}
