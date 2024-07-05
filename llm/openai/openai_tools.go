package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
)

type ParametersProperties struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Parameters struct {
	Type       string                          `json:"type"`
	Properties map[string]ParametersProperties `json:"properties"`
	Required   []string                        `json:"required"`
}

func (client *Client) SetTools(tools []llm.ToolDefinition) {
	client.tools = make(map[string]llm.ToolDefinition)

	for _, tool := range tools {
		client.tools[tool.Name] = tool
	}
}

func (client *Client) getTools() (tools []openai.Tool) {

	for _, tool := range client.tools {
		properties := make(map[string]ParametersProperties)
		for name, prop := range tool.Parameters.Properties {
			properties[name] = ParametersProperties{
				Type:        prop.Type,
				Description: prop.Description,
			}
		}

		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters: Parameters{
					Type:       tool.Parameters.Type,
					Properties: properties,
					Required:   tool.Parameters.Required,
				},
			},
		})
	}

	return tools
}

func (client *Client) callTool(ctx context.Context, name, arguments string) (string, error) {
	tool, ok := client.tools[name]
	if !ok {
		return "", fmt.Errorf("unknown tool %s", name)
	}

	var input map[string]interface{}
	if arguments != "" {
		err := json.Unmarshal([]byte(arguments), &input)
		if err != nil {
			return "", err
		}
	}

	return tool.Call(ctx, input)
}
