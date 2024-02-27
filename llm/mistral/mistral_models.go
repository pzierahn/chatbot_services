package mistral

import "strings"

const prefix = "mistral."

const (
	Large  = "mistral-large-latest"
	Medium = "mistral-medium-latest"
	Small  = "mistral-small-latest"
)

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
