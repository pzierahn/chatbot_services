package test

import (
	"context"
	"fmt"
	"github.com/pzierahn/brainboost/account"
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

	amount1 := 1000
	ctx, userId := setup.createRandomSignInWithFunding(amount1)
	defer setup.DeleteUser(userId)

	payments, err := setup.account.GetPayments(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	setup.Report.Run("payments_without_auth", func(t testing) bool {
		_, err := setup.account.GetPayments(context.Background(), &emptypb.Empty{})
		return t.expectError(err)
	})

	setup.Report.Run("payments_insert", func(t testing) bool {
		if len(payments.Items) != 1 {
			return t.fail(fmt.Errorf("payments_insert: payments.Items != 1"))
		}

		if payments.Items[0].Amount != uint32(amount1) {
			return t.fail(fmt.Errorf("payments_insert: payments.Items[0].Amount != %d", amount1))
		}

		return t.pass()
	})
}
