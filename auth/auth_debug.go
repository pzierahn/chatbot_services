package auth

import (
	"context"
	"github.com/google/uuid"
)

type Debug struct {
	user uuid.UUID
}

func WithUser(user uuid.UUID) Service {
	return &Debug{
		user: user,
	}
}

func (auth *Debug) ValidateToken(_ context.Context) (uuid.UUID, error) {
	return auth.user, nil
}
