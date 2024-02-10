package account

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Usage struct {
	Id     uuid.UUID
	UserId string
	Model  string
	Input  uint32
	Output uint32
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
