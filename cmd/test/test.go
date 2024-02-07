package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"net/http"
	"os"
)

const credentialsFile = "service_account.json"
const (
	verifyCustomTokenURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=%s"
)

// https://github.com/firebase/firebase-admin-go/blob/1d2a52c3c8195451b5ad2e0a173906bd6eb9529d/integration/auth/auth_test.go#L199
func postRequest(url string, req []byte) ([]byte, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func signInWithCustomToken(token string) (string, error) {
	req, err := json.Marshal(map[string]interface{}{
		"token":             token,
		"returnSecureToken": true,
	})
	if err != nil {
		return "", err
	}

	resp, err := postRequest(fmt.Sprintf(verifyCustomTokenURL, os.Getenv("CHATBOT_API_KEY")), req)
	if err != nil {
		return "", err
	}
	var respBody struct {
		IDToken string `json:"idToken"`
	}
	if err := json.Unmarshal(resp, &respBody); err != nil {
		return "", err
	}
	return respBody.IDToken, err
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	ctx := context.Background()

	var opts []option.ClientOption

	if _, err := os.Stat(credentialsFile); err == nil {
		serviceAccount := option.WithCredentialsFile(credentialsFile)
		opts = append(opts, serviceAccount)
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		log.Fatalf("failed to create firebase app: %v", err)
	}

	fireAuth, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("failed to create firebase fireAuth: %v", err)
	}

	token, err := fireAuth.CustomToken(ctx, "j7jjxLD9rla2DrZoeUu3Tnft4812")
	if err != nil {
		log.Fatalf("failed to create custom token: %v", err)
	}

	idToken, err := signInWithCustomToken(token)
	if err != nil {
		log.Fatalf("failed to sign in with custom token: %v", err)
	}

	conn, err := grpc.Dial(
		"brainboost-services-2qkjmuus4a-ey.a.run.app:443",
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() { _ = conn.Close() }()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "Bearer "+idToken)

	account := proto.NewAccountServiceClient(conn)
	costs, err := account.GetCosts(ctx, &empty.Empty{})
	if err != nil {
		log.Fatalf("could not get costs: %v", err)
	}

	log.Println(costs)
}
