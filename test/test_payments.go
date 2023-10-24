package test

import (
	"context"
	"fmt"
	"github.com/pzierahn/brainboost/account"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (setup *Setup) paymentsCreateFounds(ctx context.Context, userId string, amount int) {
	_, err := setup.db.Exec(ctx, `INSERT INTO payments (user_id, amount)
			values ($1, $2)`, userId, amount)
	if err != nil {
		log.Fatal(err)
	}
}

func (setup *Setup) paymentsCreateUsage(ctx context.Context, userId string, usage account.Usage) {
	_, err := setup.db.Exec(ctx,
		`INSERT INTO openai_usage (user_id, model, input_tokens, output_tokens)
			VALUES ($1, $2, $3, $4)
			RETURNING id`,
		userId, usage.Model, usage.Input, usage.Output)
	if err != nil {
		log.Fatal(err)
	}
}

func (setup *Setup) Payments() {

	setup.Report.Run("payments_without_auth", func(t testing) bool {
		_, err := setup.account.GetPayments(context.Background(), &emptypb.Empty{})
		return t.expectError(err)
	})

	setup.Report.Run("payments_accounting", func(t testing) bool {
		amount1 := 1000
		ctx, userId := setup.createRandomSignInWithFunding(amount1)
		defer setup.DeleteUser(userId)

		payments, err := setup.account.GetPayments(ctx, &emptypb.Empty{})
		if err != nil {
			log.Fatal(err)
		}

		if len(payments.Items) != 1 {
			return t.fail(fmt.Errorf("payments_insert: payments.Items != 1"))
		}

		if payments.Items[0].Amount != uint32(amount1) {
			return t.fail(fmt.Errorf("payments_insert: payments.Items[0].Amount != %d", amount1))
		}

		return t.pass()
	})

	setup.Report.Run("payments_with_founding", func(t testing) bool {
		amount1 := 1000
		ctx, userId := setup.createRandomSignInWithFunding(amount1)
		defer setup.DeleteUser(userId)

		coll, err := setup.collections.Create(ctx, &pb.Collection{
			Name: "Test",
		})
		if err != nil {
			return t.fail(err)
		}

		message, err := setup.chat.Chat(ctx, &pb.Prompt{
			Prompt:       "Say something hello",
			CollectionId: coll.Id,
			Options: &pb.PromptOptions{
				Model:       openai.GPT3Dot5Turbo,
				Temperature: 0,
				MaxTokens:   10,
			},
		})
		if err != nil {
			return t.fail(err)
		}

		if len(message.Text) == 0 {
			return t.fail(fmt.Errorf("message.Text == 0"))
		}

		return t.pass()
	})
}
