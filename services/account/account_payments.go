package account

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service *Service) getPayments(ctx context.Context, userId string) (*pb.Payments, error) {
	payments, err := service.Database.GetPayments(ctx, userId)
	if err != nil {
		return nil, err
	}

	pbPayments := make([]*pb.Payment, len(payments))
	for idx, pay := range payments {
		pbPayments[idx] = &pb.Payment{
			Id:     pay.Id.String(),
			Date:   timestamppb.New(pay.Date),
			Amount: uint32(pay.Amount),
		}
	}

	return &pb.Payments{
		Items: pbPayments,
	}, nil
}

func (service *Service) GetPayments(ctx context.Context, _ *emptypb.Empty) (*pb.Payments, error) {
	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	return service.getPayments(ctx, userId)
}
