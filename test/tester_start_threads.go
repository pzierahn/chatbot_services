package test

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/proto"
)

const (
	openaiGPT4TurboPreview = "openai.gpt-4-turbo-preview"
	openaiGPT3DOT5Turbo16k = "openai.gpt-3.5-turbo-16k"
	amazonTitanTextExpress = "amazon.titan-text-express-v1"
	anthropicClaudeV2      = "anthropic.claude-v2"
	googleGeminiPro        = "google.gemini-pro"
	mistral                = "mistral.mistral-large-latest"
)

var testModelOptions = &pb.ModelOptions{
	Model: googleGeminiPro,
	TopP:  1.0,
}

func (test Tester) TestStartThread() {
	test.runTest("TestThread_start", func(ctx context.Context) error {
		collection, err := test.collections.Create(ctx, &pb.Collection{
			Name: "test",
		})
		if err != nil {
			return err
		}

		thread, err := test.chat.StartThread(ctx, &pb.ThreadPrompt{
			Prompt:       "Say Hello!",
			CollectionId: collection.Id,
			ModelOptions: testModelOptions,
		})
		if err != nil {
			return err
		}

		if thread.Id == "" {
			return fmt.Errorf("thread id missing")
		}

		if thread.Messages[0].Prompt != "Say Hello!" {
			return fmt.Errorf("unexpected prompt: %v", thread.Messages[0].Prompt)
		}

		return nil
	})
}
