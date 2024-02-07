package auth

import (
	"context"
)

type Service interface {
	ValidateToken(ctx context.Context) (uid string, err error)
}
