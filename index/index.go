package index

import (
	"github.com/pzierahn/braingain/database"
	"github.com/sashabaranov/go-openai"
	storage_go "github.com/supabase-community/storage-go"
)

type Index struct {
	DB      *database.Client
	GPT     *openai.Client
	Storage *storage_go.Client
}
