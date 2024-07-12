package account

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Usage struct {
	Id     uuid.UUID
	UserId string
	Model  string
	Input  uint32
	Output uint32
}

func (service *Service) getUsage(ctx context.Context, userId string) (*pb.Usage, error) {
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

	var usage []*pb.ModelUsage

	for modelId, calls := range modelCalls {
		price := prices[modelId]
		input := inputs[modelId]
		output := outputs[modelId]

		usage = append(usage, &pb.ModelUsage{
			Model:    modelId,
			Input:    input,
			Output:   output,
			Costs:    price.Cost(input, output),
			Requests: calls,
		})
	}

	return &pb.Usage{
		Models: usage,
	}, nil
}

func (service *Service) GetCosts(ctx context.Context, _ *emptypb.Empty) (*pb.Usage, error) {
	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	return service.getUsage(ctx, userId)
}
