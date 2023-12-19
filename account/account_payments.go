package account

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (service *Service) GetPayments(ctx context.Context, _ *emptypb.Empty) (*pb.Payments, error) {
	userId, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := service.db.Query(
		ctx,
		`SELECT id, date, amount 
			FROM payments
			WHERE user_id = $1 
		  	ORDER BY date DESC`,
		userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := &pb.Payments{}

	for rows.Next() {
		payment := &pb.Payment{}
		var date time.Time

		err = rows.Scan(
			&payment.Id,
			&date,
			&payment.Amount,
		)
		if err != nil {
			return nil, err
		}
		payment.Date = timestamppb.New(date)

		payments.Items = append(payments.Items, payment)
	}

	return payments, nil
}
