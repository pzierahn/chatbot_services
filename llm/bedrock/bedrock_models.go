package bedrock

import "strings"

const (
	ClaudeV2Dot1 = "anthropic.claude-v2:1"
	ClaudeV3     = "anthropic.claude-3-sonnet-20240229-v1:0"
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
