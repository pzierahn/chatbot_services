package openai

import (
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
	"strings"
)

const modelPrefix = "openai."

const (
	LargeEmbedding3 = openai.LargeEmbedding3
	SmallEmbedding3 = openai.SmallEmbedding3
)

const (
	DimensionModelLarge = 3072
	DimensionModelSmall = 1536
)

var ModelCosts = map[string]llm.PricePer1000Tokens{
	//
	// GPT-4 Turbo
	//
	"gpt-4o-mini": {
		Input:  0.00015,
		Output: 0.00060,
	},
	"gpt-4o-2024-08-06": {
		Input:  0.0025,
		Output: 0.0100,
	},
	openai.GPT4o: {
		Input:  0.005,
		Output: 0.015,
	},
	"openai.gpt-4o": {
		Input:  0.005,
		Output: 0.015,
	},
	openai.GPT4o20240513: {
		Input:  0.005,
		Output: 0.015,
	},
	openai.GPT4Turbo0125: {
		Input:  0.01,
		Output: 0.03,
	},
	openai.GPT4Turbo1106: {
		Input:  0.01,
		Output: 0.03,
	},
	"gpt-4-1106-vision-preview": {
		Input:  0.01,
		Output: 0.03,
	},
	openai.GPT4VisionPreview: {
		Input:  0.01,
		Output: 0.03,
	},
	"gpt-4-turbo": {
		Input:  10.00 / 1_000,
		Output: 30.00 / 1_000,
	},
	"gpt-4-turbo-2024-04-09": {
		Input:  10.00 / 1_000,
		Output: 30.00 / 1_000,
	},
	//
	// GPT-4
	//
	openai.GPT4: {
		Input:  0.03,
		Output: 0.06,
	},
	openai.GPT40613: {
		Input:  0.03,
		Output: 0.06,
	},
	openai.GPT432K: {
		Input:  0.06,
		Output: 0.12,
	},
	//
	// GPT-3.5 Turbo
	//
	openai.GPT3Dot5Turbo0125: {
		Input:  0.0005,
		Output: 0.0015,
	},
	openai.GPT3Dot5Turbo1106: {
		Input:  0.0010,
		Output: 0.0020,
	},
	openai.GPT3Dot5Turbo0613: {
		Input:  0.0015,
		Output: 0.0020,
	},
	openai.GPT3Dot5Turbo16K: {
		Input:  0.0030,
		Output: 0.0040,
	},
	openai.GPT3Dot5Turbo16K0613: {
		Input:  0.0030,
		Output: 0.0040,
	},
	openai.GPT3Dot5Turbo0301: {
		Input:  0.0015,
		Output: 0.0020,
	},
	//
	// Embedding models
	//
	string(openai.AdaEmbeddingV2): {
		Input:  0.00002,
		Output: 0.0,
	},
	string(LargeEmbedding3): {
		Input:  0.00013,
		Output: 0.0,
	},
}

func (client *Client) ProvidesModel(name string) bool {
	_, ok := ModelCosts[name]

	switch {
	case strings.HasPrefix(name, modelPrefix):
		return true
	case ok:
		return true
	default:
		return false
	}
}
