package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	supa "github.com/nedpals/supabase-go"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"time"
)

type AdminCreateUserRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	EmailConfirm bool   `json:"email_confirm"`
}

type AdminCreateUserResponse struct {
	Id               string    `json:"id"`
	Aud              string    `json:"aud"`
	Role             string    `json:"role"`
	Email            string    `json:"email"`
	EmailConfirmedAt time.Time `json:"email_confirmed_at"`
	Phone            string    `json:"phone"`
	AppMetadata      struct {
		Provider  string   `json:"provider"`
		Providers []string `json:"providers"`
	} `json:"app_metadata"`
	UserMetadata struct {
	} `json:"user_metadata"`
	Identities []struct {
		Id           string `json:"id"`
		UserId       string `json:"user_id"`
		IdentityData struct {
			Email string `json:"email"`
			Sub   string `json:"sub"`
		} `json:"identity_data"`
		Provider     string    `json:"provider"`
		LastSignInAt time.Time `json:"last_sign_in_at"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	} `json:"identities"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Credentials struct {
	Id       string
	Email    string
	Password string
}

func (setup *Setup) CreateUser() (user Credentials) {
	url := setup.SupabaseUrl + "/auth/v1/admin/users"

	userName := fmt.Sprintf("user-%x", uuid.New().ID())

	reqBody := AdminCreateUserRequest{
		Email:        userName + "@example.com",
		Password:     uuid.NewString(),
		EmailConfirm: true,
	}

	byt, _ := json.Marshal(reqBody)

	// Post http request to create user
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(byt))
	if err != nil {
		log.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+setup.Token)
	req.Header.Set("apikey", setup.Token)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var response AdminCreateUserResponse
	// Read response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	return Credentials{
		Id:       response.Id,
		Email:    response.Email,
		Password: reqBody.Password,
	}
}

func (setup *Setup) DeleteUser(id string) {
	url := setup.SupabaseUrl + "/auth/v1/admin/users/" + id

	reqBody := struct{}{}
	byt, _ := json.Marshal(reqBody)

	// Post http request to create user
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(byt))
	if err != nil {
		log.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+setup.Token)
	req.Header.Set("apikey", setup.Token)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Delete user %s: %s", id, resp.Status)
}

func (setup *Setup) createRandomSignIn() (context.Context, string) {
	user := setup.CreateUser()

	supabase := supa.CreateClient(setup.SupabaseUrl, setup.Token)
	details, err := supabase.Auth.SignIn(context.Background(), supa.UserCredentials{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		"Authorization": "Bearer " + details.AccessToken,
	}))

	return ctx, user.Id
}
