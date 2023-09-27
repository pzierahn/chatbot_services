package auth

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	ValidateToken(ctx context.Context) (uuid.UUID, error)
}
