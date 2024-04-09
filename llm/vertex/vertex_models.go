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
		Input:  0.50 / 1_000_000,
		Output: 1.50 / 1_000_000,
	},
	GeminiPro15: {
		Input:  7 / 1_000_000,
		Output: 21 / 1_000_000,
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
