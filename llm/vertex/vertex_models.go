package vertex

import (
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

const modelPrefix = "google."

const (
	GeminiPro15        = "gemini-1.5-pro-002"
	GeminiFlash        = "gemini-1.5-flash-002"
	GeminiExperimental = "gemini-experimental"
)

var ModelCosts = map[string]llm.PricePer1000Tokens{
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
	case name == GeminiPro15:
		return true
	default:
		return false
	}
}
