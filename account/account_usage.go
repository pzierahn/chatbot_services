package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

type Usage struct {
	Id     uuid.UUID
	UserId string
	Model  string
	Input  uint32
	Output uint32
}

func (service *Service) Track(ctx context.Context, usage llm.ModelUsage) {
	if usage.UserId == "" || usage.InputTokens == 0 {
		return
	}

	_, err := service.db.Exec(
		ctx,
		`INSERT INTO model_usages (user_id, model, input_tokens, output_tokens) 
			VALUES ($1, $2, $3, $4)`,
		usage.UserId,
		usage.Model,
		usage.InputTokens,
		usage.OutputTokens,
	)
	if err != nil {
		log.Printf("failed to record usage: %v", err)
	}
}

func (service *Service) GetCosts(ctx context.Context, _ *emptypb.Empty) (*pb.Costs, error) {
	userID, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := service.db.Query(ctx,
		`SELECT model, COUNT(model), SUM(input_tokens), SUM(output_tokens)
			FROM model_usages
			WHERE user_id = $1
			GROUP BY model`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var costs pb.Costs
	for rows.Next() {
		var model pb.ModelCosts
		err = rows.Scan(
			&model.Model,
			&model.Requests,
			&model.Input,
			&model.Output,
		)
		if err != nil {
			return nil, err
		}

		modelPrice := prices[model.Model]
		model.Costs = modelPrice.cost(model.Input, model.Output)

		costs.Models = append(costs.Models, &model)
	}

	return &costs, nil
}
