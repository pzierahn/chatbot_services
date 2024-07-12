package llm

// ModelUsage defines the usage of a model
type ModelUsage struct {
	// Model name
	Model string `json:"model,omitempty"`

	// UserId is the user id
	UserId string `json:"user_id,omitempty"`

	// InputTokens is the number of tokens in the prompt
	InputTokens uint32 `json:"prompt_tokens,omitempty"`

	// OutputTokens is the number of tokens in the completion
	OutputTokens uint32 `json:"completion_tokens,omitempty"`
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
