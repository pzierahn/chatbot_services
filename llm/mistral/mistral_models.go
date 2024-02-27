package mistral

import "strings"

const prefix = "mistral."

const (
	MistralLarge  = prefix + "mistral-large-latest"
	MistralMedium = prefix + "mistral-medium-latest"
	MistralSmall  = prefix + "mistral-small-latest"
)

func (client *Client) ProvidesModel(name string) bool {
	switch {
	case strings.HasPrefix(name, prefix):
		return true
	default:
		return false
	}
}
