package account

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (service *Service) getFinancialSummary(ctx context.Context, userId string) (*pb.FinancialSummary, error) {
	payments, err := service.getPayments(ctx, userId)
	if err != nil {
		return nil, err
	}

	costs, err := service.getUsage(ctx, userId)
	if err != nil {
		return nil, err
	}

	overview := &pb.FinancialSummary{
		Payments: payments.Items,
		Usage:    costs.Models,
	}

	for _, payment := range payments.Items {
		overview.Balance += int32(payment.Amount)
	}

	for _, model := range costs.Models {
		overview.Balance -= int32(model.Costs)
	}

	return overview, nil
}

func (service *Service) GetFinancialSummary(ctx context.Context, _ *emptypb.Empty) (*pb.FinancialSummary, error) {
	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	return service.getFinancialSummary(ctx, userId)
}
