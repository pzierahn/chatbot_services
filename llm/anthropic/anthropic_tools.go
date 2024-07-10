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

type toolConverter []*llm.ToolDefinition

// toClaude converts a list of tool definitions to a list of Claude tools.
func (list toolConverter) toClaude() []ClaudeTool {
	tools := make([]ClaudeTool, len(list))

	for inx, tool := range list {
		properties := make(map[string]ClaudeToolProperty)
		for name, prop := range tool.Parameters.Properties {
			properties[name] = ClaudeToolProperty{
				Type:        prop.Type,
				Description: prop.Description,
			}
		}

		tools[inx] = ClaudeTool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: ClaudeToolInput{
				Type:       tool.Parameters.Type,
				Properties: properties,
				Required:   tool.Parameters.Required,
			},
		}
	}

	return tools
}

// getFunction returns the function call for a tool by name.
func (list toolConverter) getFunction(name string) (llm.FunctionCall, bool) {
	for _, tool := range list {
		if tool.Name == name {
			return tool.Call, true
		}
	}

	return nil, false
}

func getToolConfig(config *llm.ToolChoice) *ToolChoice {
	if config == nil {
		return nil
	}

	switch config.Type {
	case llm.ToolUseNone:
		// Not supported by Claude.
		return nil
	default:
		return &ToolChoice{
			Type: config.Type,
			Name: config.Name,
		}
	}
}
