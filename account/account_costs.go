package account

import "github.com/sashabaranov/go-openai"

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
}
