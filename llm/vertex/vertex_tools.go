package vertex

import (
	"cloud.google.com/go/vertexai/genai"
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

type toolConverter []*llm.ToolDefinition

func (list toolConverter) toVertex() (tools []*genai.Tool) {
	var functions []*genai.FunctionDeclaration

	for _, tool := range list {
		properties := make(map[string]*genai.Schema)
		for name, prop := range tool.Parameters.Properties {
			properties[name] = &genai.Schema{
				Type:        genai.TypeString,
				Description: prop.Description,
			}
		}

		functions = append(functions, &genai.FunctionDeclaration{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters: &genai.Schema{
				Type:       genai.TypeObject,
				Properties: properties,
				Required:   tool.Parameters.Required,
			},
		})
	}

	return []*genai.Tool{{
		FunctionDeclarations: functions,
	}}
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
