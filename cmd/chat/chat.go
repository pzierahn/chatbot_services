package main

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/llm/bedrock"
	"log"
)

// main demonstrates how to use the bedrock client to generate completions.
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := bedrock.New(llm.DummyTracker{})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	resp, err := client.GenerateCompletion(ctx, &llm.GenerateRequest{
		Messages: []*llm.Message{{
			Type: llm.MessageTypeUser,
			Text: "What is the meaning of life?",
		}},
		Model:       bedrock.ClaudeSonnet,
		MaxTokens:   128,
		TopP:        1.0,
		Temperature: 1.0,
	})
	if err != nil {
		log.Fatalf("Failed to generate completion: %v", err)
	}

	log.Printf("Response: %v", resp.Text)
}
