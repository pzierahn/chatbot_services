package auth

import (
	"context"
)

type insecureService struct {
	uid string
}

func (service insecureService) ValidateToken(_ context.Context) (uid string, err error) {
	// Allow all
	return service.uid, nil
}

func WithUser(uid string) (service Service, err error) {
	return &insecureService{uid}, nil
}
