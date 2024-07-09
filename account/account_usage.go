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

func (service *LiveService) getCosts(ctx context.Context, userId string) (*pb.Costs, error) {
	usages, err := service.Database.GetModelUsages(ctx, userId)
	if err != nil {
		return nil, err
	}

	modelCalls := make(map[string]uint32)
	inputs := make(map[string]uint32)
	outputs := make(map[string]uint32)

	for _, usage := range usages {
		modelCalls[usage.ModelId]++
		inputs[usage.ModelId] += usage.InputTokens
		outputs[usage.ModelId] += usage.OutputTokens
	}

	var costs []*pb.ModelCosts

	for modelId, calls := range modelCalls {
		price := prices[modelId]
		input := inputs[modelId]
		output := outputs[modelId]

		costs = append(costs, &pb.ModelCosts{
			Model:    modelId,
			Input:    input,
			Output:   output,
			Costs:    price.Cost(input, output),
			Requests: calls,
		})
	}

	return &pb.Costs{
		Models: costs,
	}, nil
}

func (service *LiveService) GetCosts(ctx context.Context, _ *emptypb.Empty) (*pb.Costs, error) {
	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	return service.getCosts(ctx, userId)
}
