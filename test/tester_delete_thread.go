package test

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/services/proto"
)

func (test Tester) TestDeleteThread() {
	test.runTest("TestDeleteThread", func(ctx context.Context) error {
		collection, err := test.collections.Create(ctx, &pb.Collection{
			Name: "test",
		})
		if err != nil {
			return err
		}

		prompt := &pb.ThreadPrompt{
			Prompt:       "Say Hello",
			CollectionId: collection.Id,
			ModelOptions: testModelOptions,
		}

		thread1, err := test.chat.StartThread(ctx, prompt)
		if err != nil {
			return err
		}

		thread2, err := test.chat.StartThread(ctx, prompt)
		if err != nil {
			return err
		}

		_, err = test.chat.DeleteThread(ctx, &pb.ThreadID{Id: thread1.Id})
		if err != nil {
			return err
		}

		out, err := test.chat.ListThreadIDs(ctx, &pb.Collection{Id: collection.Id})
		if err != nil {
			return err
		}

		if len(out.Ids) != 1 {
			return fmt.Errorf("unexpected thread count")
		}

		if out.Ids[0] != thread2.Id {
			return fmt.Errorf("unexpected thread id")
		}

		return nil
	})
}
