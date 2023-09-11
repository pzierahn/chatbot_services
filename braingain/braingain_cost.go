package braingain

import "github.com/sashabaranov/go-openai"

type Costs struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	PromptCosts      float32
	CompletionCosts  float32
	TotalCosts       float32
}

func (chat Chat) calculateCosts(model string, usage openai.Usage) Costs {
	costs := Costs{
		PromptTokens:     usage.PromptTokens,
		CompletionTokens: usage.CompletionTokens,
		TotalTokens:      usage.TotalTokens,
	}

	inputTokens := float32(usage.PromptTokens) / float32(1000)
	outputTokens := float32(usage.CompletionTokens) / float32(1000)

	switch model {
	case openai.GPT3Dot5Turbo:
		costs.PromptCosts = inputTokens * 0.0015
		costs.CompletionCosts = outputTokens * 0.002
		costs.TotalCosts = costs.PromptCosts + costs.CompletionCosts
	case openai.GPT3Dot5Turbo16K:
		costs.PromptCosts = inputTokens * 0.003
		costs.CompletionCosts = outputTokens * 0.004
		costs.TotalCosts = costs.PromptCosts + costs.CompletionCosts
	case openai.GPT4:
		costs.PromptCosts = inputTokens * 0.03
		costs.CompletionCosts = outputTokens * 0.06
		costs.TotalCosts = costs.PromptCosts + costs.CompletionCosts
	}

	return costs
}
