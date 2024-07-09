package account

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (service *LiveService) getBalanceSheet(ctx context.Context, userId string) (*pb.BalanceSheet, error) {
	payments, err := service.getPayments(ctx, userId)
	if err != nil {
		return nil, err
	}

	costs, err := service.getCosts(ctx, userId)
	if err != nil {
		return nil, err
	}

	balanceSheet := &pb.BalanceSheet{
		Payments: payments.Items,
		Costs:    costs.Models,
	}

	for _, payment := range payments.Items {
		balanceSheet.Balance += int32(payment.Amount)
	}

	for _, model := range costs.Models {
		balanceSheet.Balance -= int32(model.Costs)
	}

	return balanceSheet, nil
}

func (service *LiveService) GetBalanceSheet(ctx context.Context, req *emptypb.Empty) (*pb.BalanceSheet, error) {
	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	return service.getBalanceSheet(ctx, userId)
}
