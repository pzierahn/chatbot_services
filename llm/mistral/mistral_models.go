package mistral

import (
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

const prefix = "mistral."

const (
	Large  = "mistral-large-latest"
	Medium = "mistral-medium-latest"
	Small  = "mistral-small-latest"
)

var ModelCosts = map[string]llm.PricePer1000Tokens{
	Large: {
		Input:  8.0 / 1_000.0,
		Output: 24.0 / 1_000.0,
	},
	Medium: {
		Input:  2.7 / 1_000.0,
		Output: 8.1 / 1_000.0,
	},
	Small: {
		Input:  2.0 / 1_000.0,
		Output: 6.0 / 1_000.0,
	},
}

func (client *Client) ProvidesModel(name string) bool {
	switch {
	case strings.HasPrefix(name, prefix):
		return true
	case name == Large:
		return true
	case name == Medium:
		return true
	case name == Small:
		return true
	default:
		return false
	}
}
