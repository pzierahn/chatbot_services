package account

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
)

type Mock struct {
}

func (mock *Mock) GetModelUsages(_ context.Context, _ *emptypb.Empty) (*pb.ModelUsages, error) {
	items := make([]*pb.ModelUsages_Usage, 0)
	for i := 0; i < 10; i++ {
		items = append(items, &pb.ModelUsages_Usage{
			Model:  fmt.Sprintf("model-%d", i),
			Input:  uint32(rand.Intn(10000)),
			Output: uint32(rand.Intn(10000)),
		})
	}

	return &pb.ModelUsages{Items: items}, nil
}

func (mock *Mock) CreateUsage(_ context.Context, _ Usage) (uuid.UUID, error) {
	return uuid.New(), nil
}

func NewMock() Service {
	return &Mock{}
}
