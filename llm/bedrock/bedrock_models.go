package bedrock

import "strings"

const (
	ClaudeV2     = "anthropic.claude-v2"
	ClaudeV2Dot1 = "anthropic.claude-v2:1"
)

func (client *Client) ProvidesModel(name string) bool {
	switch {
	case strings.HasPrefix(name, "anthropic."):
		return true
	case strings.HasPrefix(name, "amazon."):
		return true
	default:
		return false
	}
}
