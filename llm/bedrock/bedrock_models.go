package bedrock

import (
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

const (
	ClaudeV2Dot1 = "anthropic.claude-v2:1"
	ClaudeSonnet = "anthropic.claude-3-sonnet-20240229-v1:0"
	ClaudeHaiku  = "anthropic.claude-3-haiku-20240307-v1:0"
)

var ModelCosts = map[string]llm.PricePer1000Tokens{
	ClaudeV2Dot1: {
		Input:  0.008,
		Output: 0.024,
	},
	ClaudeSonnet: {
		Input:  0.003,
		Output: 0.015,
	},
	ClaudeHaiku: {
		Input:  0.00025,
		Output: 0.00125,
	},
}

func (client *Client) ProvidesModel(name string) bool {
	switch {
	case name == ClaudeV2Dot1:
		return true
	case name == ClaudeSonnet:
		return true
	case name == ClaudeHaiku:
		return true
	case strings.HasPrefix(name, "amazon."):
		return true
	default:
		return false
	}
}
