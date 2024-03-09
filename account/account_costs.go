package account

import (
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/llm/anthropic"
	"github.com/pzierahn/chatbot_services/llm/bedrock"
	"github.com/pzierahn/chatbot_services/llm/mistral"
	"github.com/pzierahn/chatbot_services/llm/openai"
	"github.com/pzierahn/chatbot_services/llm/vertex"
)

var prices = map[string]llm.PricePer1000Tokens{}

func init() {
	for name, price := range anthropic.ModelCosts {
		prices[name] = price
	}

	for name, price := range bedrock.ModelCosts {
		prices[name] = price
	}

	for name, price := range mistral.ModelCosts {
		prices[name] = price
	}

	for name, price := range openai.ModelCosts {
		prices[name] = price
	}

	for name, price := range vertex.ModelCosts {
		prices[name] = price
	}
}
