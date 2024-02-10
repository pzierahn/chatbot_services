package openai

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sashabaranov/go-openai"
	"os"
)

type Client struct {
	client *openai.Client
	db     *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Client {
	token := os.Getenv("OPENAI_API_KEY")
	return &Client{
		client: openai.NewClient(token),
		db:     db,
	}
}
