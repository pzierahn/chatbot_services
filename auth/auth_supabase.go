package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

type SupabaseAuth struct {
	jwtSecret []byte
}

func WithSupabase(jwtSec string) Service {
	return &SupabaseAuth{
		jwtSecret: []byte(jwtSec),
	}
}

func (auth *SupabaseAuth) ValidateToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata missing")
	}

	var tokens []string

	// Fix ESPv2 Authorization override:
	// https://stackoverflow.com/questions/59925121/google-endpoints-error-firebase-id-token-has-incorrect-aud-audience-claim
	tokens = md.Get("X-Forwarded-Authorization")
	if len(tokens) == 0 {
		tokens = md.Get("Authorization")
	}

	if len(tokens) == 0 {
		return "", fmt.Errorf("authorization missing")
	}

	bearer := tokens[0][len("Bearer "):]

	var claims jwt.MapClaims
	token, err := jwt.ParseWithClaims(bearer, &claims, func(token *jwt.Token) (interface{}, error) {
		return auth.jwtSecret, nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("invalid user id: %v", err)
	}

	return subject, nil
}
