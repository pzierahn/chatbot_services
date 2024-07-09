package test

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/services/proto"
)

func (test Tester) TestDeleteThreadMessage() {
	test.runTest("TestDeleteThreadMessage", func(ctx context.Context) error {
		collection, err := test.collections.Create(ctx, &pb.Collection{
			Name: "Test",
		})
		if err != nil {
			return err
		}

		thread, err := test.chat.StartThread(ctx, &pb.ThreadPrompt{
			Prompt:       "I have a little green rectangular object in a yellow box",
			CollectionId: collection.Id,
			ModelOptions: testModelOptions,
		})
		if err != nil {
			return err
		}

		msg, err := test.chat.PostMessage(ctx, &pb.Prompt{
			Prompt:       "What is the color of the rectangular object in the yellow box?",
			ThreadID:     thread.Id,
			ModelOptions: testModelOptions,
		})
		if err != nil {
			return err
		}

		_, err = test.chat.DeleteMessageFromThread(ctx, &pb.MessageID{
			Id: msg.Id,
		})
		if err != nil {
			return err
		}

		updatedThread, err := test.chat.GetThread(ctx, &pb.ThreadID{
			Id: thread.Id,
		})
		if err != nil {
			return err
		}

		if len(updatedThread.Messages) != 1 {
			return fmt.Errorf("expected 1 message, got %d", len(updatedThread.Messages))
		}

		return nil
	})
}
