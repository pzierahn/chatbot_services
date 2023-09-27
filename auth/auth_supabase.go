package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"os"
)

type SupabaseAuth struct {
	jwtSecret []byte
}

func WithSupabase() Service {
	return &SupabaseAuth{
		jwtSecret: []byte(os.Getenv("SUPABASE_JWT_SECRET")),
	}
}

func (auth *SupabaseAuth) ValidateToken(ctx context.Context) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, fmt.Errorf("metadata missing")
	}

	var tokens []string

	// Fix ESPv2 Authorization override:
	// https://stackoverflow.com/questions/59925121/google-endpoints-error-firebase-id-token-has-incorrect-aud-audience-claim
	tokens = md.Get("X-Forwarded-Authorization")
	if len(tokens) == 0 {
		tokens = md.Get("Authorization")
	}

	if len(tokens) == 0 {
		return uuid.Nil, fmt.Errorf("authorization missing")
	}

	bearer := tokens[0][len("Bearer "):]

	var claims jwt.MapClaims
	token, err := jwt.ParseWithClaims(bearer, &claims, func(token *jwt.Token) (interface{}, error) {
		return auth.jwtSecret, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id: %v", err)
	}

	id, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id: %v", err)
	}

	return id, nil
}
