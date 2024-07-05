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

	params := Parameters{
		Type: "object",
		Properties: map[string]ParametersProperties{
			"prompt": {
				Type:        "string",
				Description: "The prompt for which to retrieve sources.",
			},
		},
		Required: []string{"prompt"},
	}
	byt, _ := json.Marshal(params)

	thread := []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleSystem,
		Content: "You are a helpful assistant. Quote the sources by \\cite{SourceID}",
	}, {
		Role:    openai.ChatMessageRoleUser,
		Content: "Who is Arnold Pitterson?",
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
				Description: "Retrieves the sources for the prompt.",
				Parameters:  json.RawMessage(byt),
			},
		}},
	})
	if err != nil {
		log.Fatal(err)
	}

	byt, _ = json.MarshalIndent(resp, "", "  ")
	log.Println(string(byt))

	sources := []map[string]string{{
		"SourceID": "2131245",
		"Content":  "Arnold Pitterson is a fictional character in the book 'The City of Glass' by Paul Auster.",
	}}

	if resp.Choices[0].FinishReason == openai.FinishReasonToolCalls {
		thread = append(thread, resp.Choices[0].Message)
		tool := resp.Choices[0].Message.ToolCalls[0]
		function := tool.Function

		// Add your code here to handle the function call
		log.Println("Function call:", function.Name)
		log.Println("Function parameters:", function.Arguments)

		sourceByt, _ := json.Marshal(sources)
		thread = append(thread, openai.ChatCompletionMessage{
			ToolCallID: tool.ID,
			Role:       openai.ChatMessageRoleTool,
			Content:    string(sourceByt),
		})

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
