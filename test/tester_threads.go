package test

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/proto"
	"log"
)

func (test Tester) TestThreads() {
	test.runTest("TestThreads_start", func(ctx context.Context) error {
		collection, err := test.collections.Create(ctx, &pb.Collection{
			Name: "test",
		})
		if err != nil {
			return err
		}

		thread, err := test.chat.StartThread(ctx, &pb.ThreadPrompt{
			Prompt:       "Say Hello",
			CollectionId: collection.Id,
			ModelOptions: &pb.ModelOptions{
				Model: "gemini-pro",
			},
			Threshold: 0.2,
			Limit:     1,
		})
		if err != nil {
			return err
		}

		if thread.Completion == "" {
			return fmt.Errorf("completion missing")
		}

		return nil
	})

	test.runTest("TestThreads_message", func(ctx context.Context) error {
		collection, err := test.collections.Create(ctx, &pb.Collection{
			Name: "Test",
		})
		if err != nil {
			return err
		}

		thread, err := test.chat.StartThread(ctx, &pb.ThreadPrompt{
			Prompt:       "I have a little green object in a yellow box",
			CollectionId: collection.Id,
			ModelOptions: &pb.ModelOptions{
				Model: "gemini-pro",
			},
			Threshold: 0.2,
			Limit:     1,
		})
		if err != nil {
			return err
		}

		log.Printf("Thread: %+v", thread)
		message, err := test.chat.PostMessage(ctx, &pb.Prompt{
			Prompt:   "What color is the object?",
			ThreadID: thread.Id,
			ModelOptions: &pb.ModelOptions{
				Model: "gemini-pro",
			},
		})
		if err != nil {
			return err
		}

		log.Printf("Message: %+v", message)

		return nil
	})
}
