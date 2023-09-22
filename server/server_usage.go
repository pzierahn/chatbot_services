package server

import (
	"context"
	"github.com/pzierahn/brainboost/auth"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) GetModelUsages(ctx context.Context, _ *emptypb.Empty) (*pb.ModelUsages, error) {
	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	usages, err := server.db.GetModelUsages(ctx, uid)
	if err != nil {
		return nil, err
	}

	data := &pb.ModelUsages{}
	for _, usage := range usages {
		data.Items = append(data.Items, &pb.ModelUsages_Usage{
			Model:  usage.Model,
			Input:  usage.Input,
			Output: usage.Output,
		})
	}

	return data, nil
}
