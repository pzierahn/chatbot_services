package index

import (
	"github.com/pzierahn/braingain/database"
	"github.com/sashabaranov/go-openai"
	"os"
)

type Index struct {
	collection string
	conn       *database.Client
	ai         *openai.Client
}

func NewIndex(conn *database.Client, collection string) *Index {
	token := os.Getenv("OPENAI_API_KEY")
	ai := openai.NewClient(token)

	return &Index{
		collection: collection,
		conn:       conn,
		ai:         ai,
	}
}
