package account

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (service *Service) getOverview(ctx context.Context, userId string) (*pb.Overview, error) {
	payments, err := service.getPayments(ctx, userId)
	if err != nil {
		return nil, err
	}

	costs, err := service.getUsage(ctx, userId)
	if err != nil {
		return nil, err
	}

	overview := &pb.Overview{
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

func (service *Service) GetOverview(ctx context.Context, _ *emptypb.Empty) (*pb.Overview, error) {
	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	return service.getOverview(ctx, userId)
}
