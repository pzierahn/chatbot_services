package anthropic

import (
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
)

type ClaudeToolProperty struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

type ClaudeToolInput struct {
	Type       string                        `json:"type,omitempty"`
	Properties map[string]ClaudeToolProperty `json:"properties,omitempty"`
	Required   []string                      `json:"required,omitempty"`
}

type ClaudeTool struct {
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	InputSchema ClaudeToolInput `json:"input_schema,omitempty"`
}

func (client *Client) SetTools(tools []llm.ToolDefinition) {
	client.tools = make(map[string]llm.ToolDefinition)

	for _, tool := range tools {
		client.tools[tool.Name] = tool
	}
}

func (client *Client) getTools() (tools []ClaudeTool) {

	for _, tool := range client.tools {
		properties := make(map[string]ClaudeToolProperty)
		for name, prop := range tool.Parameters.Properties {
			properties[name] = ClaudeToolProperty{
				Type:        prop.Type,
				Description: prop.Description,
			}
		}

		tools = append(tools, ClaudeTool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: ClaudeToolInput{
				Type:       tool.Parameters.Type,
				Properties: properties,
				Required:   tool.Parameters.Required,
			},
		})
	}

	return tools
}

func (client *Client) callTool(ctx context.Context, name string, arguments map[string]interface{}) (string, error) {
	tool, ok := client.tools[name]
	if !ok {
		return "", fmt.Errorf("unknown tool %s", name)
	}

	return tool.Call(ctx, arguments)
}
