package test

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/proto"
	"reflect"
)

func (test Tester) TestGetThread() {
	test.runTest("TestGetThread", func(ctx context.Context) error {
		collection, err := test.collections.Create(ctx, &pb.Collection{
			Name: "test",
		})
		if err != nil {
			return err
		}

		ori, err := test.chat.StartThread(ctx, &pb.ThreadPrompt{
			Prompt:       "Say Hello",
			CollectionId: collection.Id,
			ModelOptions: &pb.ModelOptions{
				Model: testModel,
			},
		})
		if err != nil {
			return err
		}

		out, err := test.chat.GetThread(ctx, &pb.ThreadID{Id: ori.Id})
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(ori, out) {
			return fmt.Errorf("threads not equal")
		}

		return nil
	})
}
