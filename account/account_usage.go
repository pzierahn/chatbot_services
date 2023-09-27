package account

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Usage struct {
	Id     uuid.UUID
	UserId uuid.UUID
	Model  string
	Input  uint32
	Output uint32
}

func (service *Service) GetModelUsages(ctx context.Context, _ *emptypb.Empty) (*pb.ModelUsages, error) {
	userID, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := service.db.Query(ctx,
		`SELECT model, SUM(input), SUM(output)
			FROM openai_usage
			WHERE user_id = $1
			GROUP BY model`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usages pb.ModelUsages
	for rows.Next() {
		var usage pb.ModelUsages_Usage
		err = rows.Scan(
			&usage.Model,
			&usage.Input,
			&usage.Output)
		if err != nil {
			return nil, err
		}

		usages.Items = append(usages.Items, &usage)
	}

	return &usages, nil
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
		`INSERT INTO openai_usage (user_id, model, input_tokens, output_tokens)
			VALUES ($1, $2, $3, $4)
			RETURNING id`,
		userID, usage.Model, usage.Input, usage.Output).
		Scan(&id)

	return id, err
}
