package bedrock

import "strings"

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
