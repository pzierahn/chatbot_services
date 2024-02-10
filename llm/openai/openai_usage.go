package openai

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"log"
)

func (client *Client) trackUsage(ctx context.Context, usage llm.ModelUsage) {
	if usage.UserId == "" {
		return
	}

	_, err := client.db.Exec(
		ctx,
		`INSERT INTO model_usages (user_id, model, input_tokens, output_tokens) 
			VALUES ($1, $2, $3, $4)`,
		usage.UserId,
		usage.Model,
		usage.PromptTokens,
		usage.CompletionTokens,
	)
	if err != nil {
		log.Printf("Error tracking usage: %v", err)
	}
}
