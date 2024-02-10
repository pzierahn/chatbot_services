package test

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/proto"
	"strings"
)

func (test Tester) TestThreadMessages() {
	test.runTest("TestThread_messages", func(ctx context.Context) error {
		collection, err := test.collections.Create(ctx, &pb.Collection{
			Name: "Test",
		})
		if err != nil {
			return err
		}

		thread, err := test.chat.StartThread(ctx, &pb.ThreadPrompt{
			Prompt:       "I have a little green rectangular object in a yellow box",
			CollectionId: collection.Id,
			ModelOptions: &pb.ModelOptions{
				Model: "gemini-pro",
			},
		})
		if err != nil {
			return err
		}

		message, err := test.chat.PostMessage(ctx, &pb.Prompt{
			Prompt:   "What is the color of the rectangular object in the yellow box?",
			ThreadID: thread.Id,
			ModelOptions: &pb.ModelOptions{
				Model: "gemini-pro",
			},
		})
		if err != nil {
			return err
		}

		completion := strings.ToLower(message.Completion)
		if !strings.Contains(completion, "green") {
			return fmt.Errorf("unexpected completion: %v", completion)
		}

		return nil
	})
}
