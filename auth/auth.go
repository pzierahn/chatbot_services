package auth

import (
	"context"
)

type Service interface {
	ValidateToken(ctx context.Context) (string, error)
}
