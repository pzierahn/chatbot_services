package account

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (service *Service) GetBalanceSheet(ctx context.Context, req *emptypb.Empty) (*pb.BalanceSheet, error) {
	payments, err := service.GetPayments(ctx, req)
	if err != nil {
		return nil, err
	}

	costs, err := service.GetCosts(ctx, req)
	if err != nil {
		return nil, err
	}

	balanceSheet := &pb.BalanceSheet{
		Payments: payments.Items,
		Costs:    costs.Models,
	}

	for _, payment := range payments.Items {
		balanceSheet.Balance += payment.Amount
	}

	for _, model := range costs.Models {
		balanceSheet.Balance -= model.Costs
	}

	return balanceSheet, nil
}

func (service *Service) HasFunding(ctx context.Context) (bool, error) {
	payments, err := service.GetPayments(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	var payedAmount uint32
	for _, payment := range payments.Items {
		payedAmount += payment.Amount
	}

	usages, err := service.GetCosts(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	var costs uint32
	for _, model := range usages.Models {
		costs += model.Costs
	}

	return payedAmount > costs, nil
}
