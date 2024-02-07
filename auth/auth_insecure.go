package auth

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
)

type insecureService struct{}

func (service insecureService) Verify(ctx context.Context) (uid string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata missing")
	}

	uids := md.Get("User-Id")
	if len(uids) != 1 {
		return "", fmt.Errorf("uid missing")
	}

	return uids[0], nil
}

func WithInsecure() (service Service, err error) {
	return &insecureService{}, nil
}
