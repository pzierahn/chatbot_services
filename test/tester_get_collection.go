package test

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/proto"
	"reflect"
)

func (test Tester) TestGetCollection() {
	test.runTest("TestGetCollection", func(ctx context.Context) error {
		_, err := test.collections.Create(ctx, &pb.Collection{
			Name: "aaaa test 2",
		})
		if err != nil {
			return err
		}

		in, err := test.collections.Create(ctx, &pb.Collection{
			Name: "test",
		})
		if err != nil {
			return err
		}

		out, err := test.collections.Get(ctx, &pb.CollectionID{
			Id: in.Id,
		})
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(in, out) {
			return fmt.Errorf("expected %v, got %v", in, out)
		}

		return nil
	})
}
