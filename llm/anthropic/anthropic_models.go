package anthropic

import "github.com/pzierahn/chatbot_services/llm"

const prefix = "anthropic."

const (
	Opus   = "claude-3-opus-20240229"
	Sonnet = "claude-3-sonnet-20240229"
)

var ModelCosts = map[string]llm.PricePer1000Tokens{
	Opus: {
		Input:  0.015,
		Output: 0.075,
	},
	Sonnet: {
		Input:  0.003,
		Output: 0.015,
	},
}

func (client *Client) ProvidesModel(name string) bool {
	switch {
	case name == Opus:
		return true
	case prefix+Opus == name:
		return true
	default:
		return false
	}
}
