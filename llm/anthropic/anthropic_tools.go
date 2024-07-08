package anthropic

import (
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

// transformToTool transforms a tool definition to a ClaudeTool object.
func transformToTool(tool llm.ToolDefinition) ClaudeTool {
	properties := make(map[string]ClaudeToolProperty)
	for name, prop := range tool.Parameters.Properties {
		properties[name] = ClaudeToolProperty{
			Type:        prop.Type,
			Description: prop.Description,
		}
	}

	return ClaudeTool{
		Name:        tool.Name,
		Description: tool.Description,
		InputSchema: ClaudeToolInput{
			Type:       tool.Parameters.Type,
			Properties: properties,
			Required:   tool.Parameters.Required,
		},
	}
}

// transformTools transforms a list of tool definitions to a map of tool names to ClaudeTool objects.
func transformTools(items []llm.ToolDefinition) []ClaudeTool {
	tools := make([]ClaudeTool, len(items))

	for inx, tool := range items {
		tools[inx] = transformToTool(tool)
	}

	return tools
}
