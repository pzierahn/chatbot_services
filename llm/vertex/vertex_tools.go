package vertex

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
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

func (client *Client) getTools() (tools []*genai.Tool) {
	for _, tool := range client.tools {
		properties := make(map[string]*genai.Schema)
		for name, prop := range tool.Parameters.Properties {
			properties[name] = &genai.Schema{
				Type:        genai.TypeString,
				Description: prop.Description,
			}
		}

		tools = append(tools, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters: &genai.Schema{
					Type:       genai.TypeObject,
					Properties: properties,
					Required:   tool.Parameters.Required,
				},
			}},
		})
	}

	return tools
}

func (client *Client) callTool(ctx context.Context, name string, arguments map[string]any) (string, error) {
	tool, ok := client.tools[name]
	if !ok {
		return "", fmt.Errorf("unknown tool %s", name)
	}

	input := make(map[string]interface{})
	for key, value := range arguments {
		input[key] = value
	}

	return tool.Call(ctx, input)
}
