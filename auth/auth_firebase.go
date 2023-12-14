package auth

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"google.golang.org/grpc/metadata"
	"strings"
)

type firebaseService struct {
	client *auth.Client
}

const credentials = "serviceAccount.json"

func WithFirebase(ctx context.Context, app *firebase.App) (service Service, err error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &firebaseService{client: client}, nil
}

func (auth *firebaseService) ValidateToken(ctx context.Context) (string, error) {
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

	bearer := strings.TrimPrefix(tokens[0], "Bearer ")
	token, err := auth.client.VerifyIDToken(ctx, bearer)
	if err != nil {
		return "", err
	}

	return token.UID, nil
}
