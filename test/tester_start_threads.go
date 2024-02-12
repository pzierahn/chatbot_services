package test

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/proto"
)

const testModel = "openai.gpt-3.5-turbo-16k"

// const testModel = "openai.gpt-4-turbo-preview"
// const testModel = "amazon.titan-text-express-v1"
// const testModel = "anthropic.claude-v2"
// const testModel = "anthropic.claude-v2"
// const testModel = "google.gemini-pro"

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
			ModelOptions: &pb.ModelOptions{
				Model: testModel,
			},
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
