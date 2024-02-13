package openai

import "github.com/sashabaranov/go-openai"

const modelPrefix = "openai."

func (client *Client) ProvidesModel(name string) bool {
	switch name {
	case modelPrefix + openai.GPT4TurboPreview:
		return true
	case modelPrefix + openai.GPT3Dot5Turbo16K:
		return true
	case openai.GPT3Dot5Turbo16K:
		return true
	case openai.GPT4TurboPreview:
		return true
	default:
		return false
	}
}
