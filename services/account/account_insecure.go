package account

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
)

// InsecureVerifier is an insecure implementation of the Verifier interface. It does not
// perform any verification, instead it trusts the user ID provided in the metadata.
type InsecureVerifier struct{}

// Verify checks if the context contains a user ID and returns it
func (verifier InsecureVerifier) Verify(ctx context.Context) (userId string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata missing")
	}

	uids := md.Get("User-Id")
	if len(uids) != 1 {
		return "", fmt.Errorf("User-Id missing")
	}

	return uids[0], nil
}

// VerifyFunding checks if the context contains a user ID and returns it. It always returns sufficient funding.
func (verifier InsecureVerifier) VerifyFunding(ctx context.Context) (userId string, err error) {
	userId, err = verifier.Verify(ctx)
	if err != nil {
		return "", err
	}

	return userId, nil
}
