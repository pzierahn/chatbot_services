package mistral

import (
	"fmt"
	"github.com/gage-technologies/mistral-go"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

type Client struct {
	client *mistral.MistralClient
	db     *pgxpool.Pool
}

func New(db *pgxpool.Pool) (*Client, error) {

	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("MISTRAL_API_KEY env var is not set")
	}

	return &Client{
		client: mistral.NewMistralClientDefault(apiKey),
		db:     db,
	}, nil
}
