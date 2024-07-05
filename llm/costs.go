package llm

type ModelUsage struct {
	Model        string `json:"model,omitempty"`
	UserId       string `json:"user_id,omitempty"`
	InputTokens  int    `json:"prompt_tokens,omitempty"`
	OutputTokens int    `json:"completion_tokens,omitempty"`
}

type PricePer1000Tokens struct {
	Input  float32
	Output float32
}

// Cost returns the cost in cents of a model usage
func (price *PricePer1000Tokens) Cost(input, output uint32) (cost uint32) {
	cost += uint32(float32(input) * price.Input)
	cost += uint32(float32(output) * price.Output)
	return cost / 10
}
