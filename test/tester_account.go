package test

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (test Tester) TestAccountCosts() {
	test.runTest("TestAccount_costs", func(ctx context.Context) error {
		collection, err := test.collections.Create(ctx, &pb.Collection{
			Name: "test",
		})
		if err != nil {
			return err
		}

		_, err = test.chat.StartThread(ctx, &pb.ThreadPrompt{
			Prompt:       "Tell a long about a pinguin",
			CollectionId: collection.Id,
			ModelOptions: &pb.ModelOptions{
				Model: testModel,
			},
		})
		if err != nil {
			return err
		}

		costs, err := test.account.GetCosts(ctx, &emptypb.Empty{})
		if err != nil {
			return err
		}

		if len(costs.Models) != 1 {
			return fmt.Errorf("expected 1 model, got %d", len(costs.Models))
		}

		if costs.Models[0].Requests != 1 {
			return fmt.Errorf("expected 1 request, got %d", costs.Models[0].Requests)
		}

		if costs.Models[0].Input == 0 {
			return fmt.Errorf("expected non-zero input, got %d", costs.Models[0].Input)
		}

		if costs.Models[0].Output == 0 {
			return fmt.Errorf("expected non-zero output, got %d", costs.Models[0].Output)
		}

		if costs.Models[0].Costs == 0 {
			return fmt.Errorf("expected non-zero costs, got %d", costs.Models[0].Costs)
		}

		return nil
	})
}
