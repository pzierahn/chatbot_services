package openai

import (
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

type toolConverter []llm.ToolDefinition

func (tools toolConverter) toOpenAI() []openai.Tool {
	items := make([]openai.Tool, len(tools))

	for inx, tool := range tools {
		properties := make(map[string]ParametersProperties)
		for name, prop := range tool.Parameters.Properties {
			properties[name] = ParametersProperties{
				Type:        prop.Type,
				Description: prop.Description,
			}
		}

		items[inx] = openai.Tool{
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
		}
	}

	return items
}

func (tools toolConverter) getFunction(name string) (llm.FunctionCall, bool) {
	for _, tool := range tools {
		if tool.Name == name {
			return tool.Call, true
		}
	}

	return nil, false
}
