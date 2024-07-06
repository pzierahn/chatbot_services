package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/llm/anthropic"
	"log"
	"strings"
)

func getSources(_ context.Context, input map[string]interface{}) (string, error) {
	prompt, ok := input["prompt"].(string)
	if !ok {
		return "", fmt.Errorf("missing prompt")
	}

	log.Printf("Call get_sources: '%s'", prompt)

	var sources []map[string]string

	if strings.Contains(prompt, "Arnold Pitterson") {
		sources = append(sources, map[string]string{
			"SourceID": "S1",
			"Content":  "Arnold Pitterson is a fictional character in the book 'The City of Glass' by Paul Auster.",
		})
	}

	if strings.Contains(prompt, "Hugo Alberts von Tahl") {
		sources = append(sources, map[string]string{
			"SourceID": "S2",
			"Content":  "Hugo Alberts von Tahl was a German philosopher who lived in the 19th century. He is known for his work on the philosophy of language and logic.",
		})
	}

	byt, _ := json.Marshal(sources)
	return string(byt), nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	client, err := anthropic.New()
	if err != nil {
		log.Fatal(err)
	}

	client.SetTools([]llm.ToolDefinition{{
		Name:        "get_sources",
		Description: "Retrieves the sources for the prompt. The prompt should be optimized for embedding retrieval. The tool will return a list of sources in JSON format with the following fields: SourceID, Content.",
		Parameters: llm.ToolParameters{
			Type: "object",
			Properties: map[string]llm.ParametersProperties{
				"prompt": {
					Type:        "string",
					Description: "The topic for which to retrieve sources. The prompt should be optimized for embedding retrieval.",
				},
			},
			Required: []string{"prompt"},
		},
		Call: getSources,
	}})

	resp, err := client.Completion(ctx, &llm.CompletionRequest{
		SystemPrompt: "You are a helpful assistant. Quote the sources by \\cite{SourceID}",
		Messages: []*llm.Message{{
			Role:    llm.MessageTypeUser,
			Content: "Who is Arnold Pitterson? Who is Hugo Alberts von Tahl?",
		}},
		Temperature: 1.0,
		TopP:        1.0,
		MaxTokens:   256,
		Model:       anthropic.ClaudeSonnet35,
	})
	if err != nil {
		log.Fatal(err)
	}

	byt, _ := json.MarshalIndent(resp, "", "  ")
	log.Println(string(byt))
}
