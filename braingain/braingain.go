package braingain

import (
	"github.com/pzierahn/braingain/database"
	"github.com/sashabaranov/go-openai"
)

type Chat struct {
	db  *database.Client
	gpt *openai.Client
}

func NewChat(db *database.Client, gpt *openai.Client) *Chat {
	return &Chat{
		db:  db,
		gpt: gpt,
	}
}
