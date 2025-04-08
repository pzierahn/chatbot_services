package vertex

import (
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

const modelPrefix = "google."

var ModelCosts = map[string]llm.PricePer1000Tokens{
	"gemini-2.0-flash": {
		Input:  0.00015,
		Output: 0.00060,
	},
	"gemini-2.0-flash-lite": {
		Input:  0.000075,
		Output: 0.00030,
	},
	"gemini-1.5-pro-002": {
		Input:  0.007,
		Output: 0.021,
	},
	"gemini-1.5-pro-001": {
		Input:  0.007,
		Output: 0.021,
	},
	"gemini-1.0-pro": {
		Input:  0.0005,
		Output: 0.0015,
	},
	"gemini-1.0-pro-001": {
		Input:  0.0005,
		Output: 0.0015,
	},
	"gemini-1.0-pro-002": {
		Input:  0.0005,
		Output: 0.0015,
	},
	"gemini-pro": {
		Input:  0.0005,
		Output: 0.0015,
	},
	"gemini-1.5-pro-preview-0409": {
		Input:  0.007,
		Output: 0.021,
	},
	"gemini-1.5-pro-preview-0514": {
		Input:  0.007,
		Output: 0.021,
	},
}

func (client *Client) ProvidesModel(name string) bool {
	switch {
	case strings.HasPrefix(name, modelPrefix):
		return true
	default:
		return false
	}
}
