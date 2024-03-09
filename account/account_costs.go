package account

import (
	"github.com/pzierahn/chatbot_services/llm/anthropic"
	"github.com/pzierahn/chatbot_services/llm/mistral"
	"github.com/sashabaranov/go-openai"
)

type price struct {
	input  float32
	output float32
}

// cost returns the cost in cents of a model usage
func (price *price) cost(input, output uint32) (cost uint32) {
	cost += uint32(float32(input) * price.input)
	cost += uint32(float32(output) * price.output)
	return cost / 10
}

var prices = map[string]price{
	//
	// GPT-4 Turbo
	//
	openai.GPT4Turbo0125: {
		input:  0.01,
		output: 0.03,
	},
	openai.GPT4Turbo1106: {
		input:  0.01,
		output: 0.03,
	},
	"gpt-4-1106-vision-preview": {
		input:  0.01,
		output: 0.03,
	},
	openai.GPT4VisionPreview: {
		input:  0.01,
		output: 0.03,
	},
	//
	// GPT-4
	//
	openai.GPT4: {
		input:  0.03,
		output: 0.06,
	},
	openai.GPT40613: {
		input:  0.03,
		output: 0.06,
	},
	openai.GPT432K: {
		input:  0.06,
		output: 0.12,
	},
	//
	// GPT-3.5 Turbo
	//
	openai.GPT3Dot5Turbo0125: {
		input:  0.0005,
		output: 0.0015,
	},
	openai.GPT3Dot5Turbo1106: {
		input:  0.0010,
		output: 0.0020,
	},
	openai.GPT3Dot5Turbo0613: {
		input:  0.0015,
		output: 0.0020,
	},
	openai.GPT3Dot5Turbo16K: {
		input:  0.0030,
		output: 0.0040,
	},
	openai.GPT3Dot5Turbo16K0613: {
		input:  0.0030,
		output: 0.0040,
	},
	openai.GPT3Dot5Turbo0301: {
		input:  0.0015,
		output: 0.0020,
	},
	//
	// Embedding models
	//
	string(openai.AdaEmbeddingV2): {
		input:  0.00002,
		output: 0.0,
	},
	string(openai.LargeEmbedding3): {
		input:  0.00013,
		output: 0.0,
	},
	//
	// Mistral AI
	//
	mistral.Large: {
		input:  8.0 / 1_000.0,
		output: 24.0 / 1_000.0,
	},
	mistral.Medium: {
		input:  2.7 / 1_000.0,
		output: 8.1 / 1_000.0,
	},
	mistral.Small: {
		input:  2.0 / 1_000.0,
		output: 6.0 / 1_000.0,
	},
	//
	// Anthropic
	//
	"claude-3-sonnet-28k-20240229": {
		input:  0.003,
		output: 0.015,
	},
	"claude-2.1": {
		input:  0.008,
		output: 0.024,
	},
	anthropic.OPUS: {
		input:  0.015,
		output: 0.075,
	},
}
