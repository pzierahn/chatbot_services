package main

import (
	"context"
	"encoding/json"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
)

type ParametersProperties struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Parameters struct {
	Type       string                          `json:"type"`
	Properties map[string]ParametersProperties `json:"properties"`
	Required   []string                        `json:"required"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	thread := []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleSystem,
		Content: "You are a helpful assistant. Quote the sources by \\cite{SourceID}",
	}, {
		Role:    openai.ChatMessageRoleUser,
		Content: "Who are Arnold Pitterson and Hugo Alberts von Tahl?",
	}}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT4o,
		MaxTokens:   1024,
		Temperature: 1.0,
		Messages:    thread,
		Tools: []openai.Tool{{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "get_sources",
				Description: "Retrieves the sources for the prompt. The prompt should be optimized for embedding retrieval. The tool will return a list of sources in JSON format with the following fields: SourceID, Content.",
				Parameters: Parameters{
					Type: "object",
					Properties: map[string]ParametersProperties{
						"prompt": {
							Type:        "string",
							Description: "The topic for which to retrieve sources. The prompt should be optimized for embedding retrieval.",
						},
					},
					Required: []string{"prompt"},
				},
			},
		}},
	})
	if err != nil {
		log.Fatal(err)
	}

	byt, _ := json.MarshalIndent(resp, "", "  ")
	log.Println(string(byt))

	sources := []map[string]string{{
		"SourceID": "S1",
		"Content":  "Arnold Pitterson is a fictional character in the book 'The City of Glass' by Paul Auster.",
	}, {
		"SourceID": "S2",
		"Content":  "Hugo Alberts von Tahl was a German philosopher who lived in the 19th century. He is known for his work on the philosophy of language and logic.",
	}}

	if resp.Choices[0].FinishReason == openai.FinishReasonToolCalls {
		thread = append(thread, resp.Choices[0].Message)

		for inx, tool := range resp.Choices[0].Message.ToolCalls {
			function := tool.Function

			// Add your code here to handle the function call
			log.Println("Function call:", function.Name)
			log.Println("Function parameters:", function.Arguments)

			sourceByt, _ := json.Marshal(sources[inx])
			thread = append(thread, openai.ChatCompletionMessage{
				ToolCallID: tool.ID,
				Role:       openai.ChatMessageRoleTool,
				Content:    string(sourceByt),
			})
		}

		resp, err = client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:       openai.GPT4o,
			MaxTokens:   1024,
			Temperature: 1.0,
			Messages:    thread,
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Final response:", resp.Choices[0].Message.Content)

		// Add the final response to thread
		thread = append(thread, resp.Choices[0].Message)

		byt, _ = json.MarshalIndent(thread, "", "  ")
		log.Println(string(byt))
	}
}
