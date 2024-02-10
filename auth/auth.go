package auth

import (
	"context"
)

type Service interface {
	Verify(ctx context.Context) (uid string, err error)
}
