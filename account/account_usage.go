package account

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/pzierahn/brainboost/proto"
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
	userID, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := service.db.Query(ctx,
		`SELECT model, SUM(input_tokens), SUM(output_tokens)
			FROM openai_usages
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
			&model.Input,
			&model.Output,
		)
		if err != nil {
			return nil, err
		}

		model.Costs += uint32(float32(model.Input)*inputCosts[model.Model]) / 10
		model.Costs += uint32(float32(model.Output)*outputCosts[model.Model]) / 10

		costs.Models = append(costs.Models, &model)
	}

	return &costs, nil
}

// CreateUsage inserts a new usage record into the openai_usage table
func (service *Service) CreateUsage(ctx context.Context, usage Usage) (uuid.UUID, error) {
	// Validate the token
	userID, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	err = service.db.QueryRow(ctx,
		`INSERT INTO openai_usages (user_id, model, input_tokens, output_tokens)
			VALUES ($1, $2, $3, $4)
			RETURNING id`,
		userID, usage.Model, usage.Input, usage.Output).
		Scan(&id)

	return id, err
}
