package vertex

import "github.com/pzierahn/chatbot_services/llm"

const modelPrefix = "google."

const (
	GeminiPro = "gemini-pro"
)

var ModelCosts = map[string]llm.PricePer1000Tokens{
	GeminiPro: {
		Input:  0.000125,
		Output: 0.000375,
	},
}

func (client *Client) ProvidesModel(name string) bool {
	switch name {
	case modelPrefix + "gemini-pro":
		return true
	default:
		return false
	}
}
