package account

import "github.com/sashabaranov/go-openai"

var inputCosts = map[string]float32{
	// GPT-3.5 Turbo
	openai.GPT3Dot5Turbo:        0.0015,
	openai.GPT3Dot5Turbo0125:    0.0015,
	openai.GPT3Dot5Turbo1106:    0.0015,
	openai.GPT3Dot5Turbo0613:    0.0015,
	openai.GPT3Dot5Turbo0301:    0.0015,
	openai.GPT3Dot5Turbo16K:     0.0015,
	openai.GPT3Dot5Turbo16K0613: 0.003,

	// GPT-4
	openai.GPT4:        0.03,
	openai.GPT40314:    0.03,
	openai.GPT40613:    0.03,
	openai.GPT432K:     0.06,
	openai.GPT432K0613: 0.06,
	openai.GPT432K0314: 0.06,

	// GPT 4 Turbo
	openai.GPT4TurboPreview:  0.01,
	openai.GPT4Turbo1106:     0.01,
	openai.GPT4Turbo0125:     0.01,
	openai.GPT4VisionPreview: 0.01,

	// Embeddings
	string(openai.AdaEmbeddingV2): 0.0001,
}

var outputCosts = map[string]float32{
	// GPT-3.5 Turbo
	openai.GPT3Dot5Turbo:        0.002,
	openai.GPT3Dot5Turbo0301:    0.002,
	openai.GPT3Dot5Turbo0613:    0.002,
	openai.GPT3Dot5Turbo16K:     0.004,
	openai.GPT3Dot5Turbo16K0613: 0.004,

	// GPT-4
	openai.GPT4:        0.06,
	openai.GPT40314:    0.06,
	openai.GPT40613:    0.06,
	openai.GPT432K:     0.12,
	openai.GPT432K0613: 0.12,
	openai.GPT432K0314: 0.12,

	// GPT 4 Turbo
	openai.GPT4TurboPreview:  0.03,
	openai.GPT4VisionPreview: 0.03,
}
