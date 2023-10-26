package account

import (
	"context"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var inputCosts = map[string]float32{
	// GPT-3.5 Turbo
	openai.GPT3Dot5Turbo:        0.0015,
	openai.GPT3Dot5Turbo0301:    0.0015,
	openai.GPT3Dot5Turbo0613:    0.0015,
	openai.GPT3Dot5Turbo16K:     0.003,
	openai.GPT3Dot5Turbo16K0613: 0.003,

	// GPT-4
	openai.GPT4:        0.03,
	openai.GPT40314:    0.03,
	openai.GPT40613:    0.03,
	openai.GPT432K:     0.06,
	openai.GPT432K0613: 0.06,
	openai.GPT432K0314: 0.06,

	// Embeddings
	openai.AdaEmbeddingV2.String(): 0.0001,
}

var outputCosts = map[string]float32{
	// GPT-3.5 Turbo
	openai.GPT3Dot5Turbo:        0.002,
	openai.GPT3Dot5Turbo0301:    0.002,
	openai.GPT3Dot5Turbo0613:    0.002,
	openai.GPT3Dot5Turbo16K:     0.004,
	openai.GPT3Dot5Turbo16K0613: 0.004,

	// GPT-4
	openai.GPT4:        0.06,
	openai.GPT40314:    0.06,
	openai.GPT40613:    0.06,
	openai.GPT432K:     0.12,
	openai.GPT432K0613: 0.12,
	openai.GPT432K0314: 0.12,
}

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
		payment := &pb.Payments_Payment{}
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

func (service *Service) HasFounding(ctx context.Context) (bool, error) {
	payments, err := service.GetPayments(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	var payedAmount uint32
	for _, payment := range payments.Items {
		payedAmount += payment.Amount
	}

	usages, err := service.GetModelUsages(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	var costs uint32
	for _, usage := range usages.Items {
		costs += uint32((float32(usage.Input) / 1000) * inputCosts[usage.Model])
		costs += uint32((float32(usage.Output) / 1000) * outputCosts[usage.Model])
	}

	return payedAmount > costs, nil
}
