package vertex

import (
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

const modelPrefix = "google."

const (
	GeminiPro   = "gemini-1.0-pro"
	GeminiPro15 = "gemini-1.5-pro-preview-0409"
)

var ModelCosts = map[string]llm.PricePer1000Tokens{
	GeminiPro: {
		Input:  0.0005,
		Output: 0.0015,
	},
	"gemini-pro": {
		Input:  0.0005,
		Output: 0.0015,
	},
	GeminiPro15: {
		Input:  0.007,
		Output: 0.021,
	},
}

func (client *Client) ProvidesModel(name string) bool {
	switch {
	case strings.HasPrefix(name, modelPrefix):
		return true
	case name == GeminiPro:
		return true
	case name == GeminiPro15:
		return true
	default:
		return false
	}
}
