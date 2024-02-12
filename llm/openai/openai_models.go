package openai

import "github.com/sashabaranov/go-openai"

func (client *Client) ProvideModel(name string) bool {
	switch name {
	case openai.GPT4TurboPreview:
		return true
	case openai.GPT3Dot5Turbo16K:
		return true
	default:
		return false
	}
}
