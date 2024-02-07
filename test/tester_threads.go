package test

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"log"
)

func (test Tester) TestThreads() {
	test.runTest("TestThreads", func(ctx context.Context) error {
		collection, err := test.collections.Create(ctx, &pb.Collection{
			Name: "test",
		})
		if err != nil {
			return err
		}

		thread, err := test.chat.StartThread(ctx, &pb.Prompt{
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

		log.Printf("Thread: %+v", thread)

		return nil
	})
}
