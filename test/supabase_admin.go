package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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

func (service Service) CreateUser() (user Credentials) {
	url := service.SupabaseUrl + "/auth/v1/admin/users"

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
	req.Header.Set("Authorization", "Bearer "+service.Token)
	req.Header.Set("apikey", service.Token)

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
