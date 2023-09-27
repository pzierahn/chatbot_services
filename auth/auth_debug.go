package auth

import (
	"context"
	"github.com/google/uuid"
)

type Debug struct {
	userId uuid.UUID
}

func WithUser(userId uuid.UUID) Service {
	return &Debug{
		userId: userId,
	}
}

func (auth *Debug) ValidateToken(_ context.Context) (uuid.UUID, error) {
	return auth.userId, nil
}
