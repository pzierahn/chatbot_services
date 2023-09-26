package main

import (
	"context"
	"fmt"
	"github.com/pzierahn/brainboost/database"
	"github.com/pzierahn/brainboost/index"
	"github.com/sashabaranov/go-openai"
	storagego "github.com/supabase-community/storage-go"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	baseDir = "./"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	ctx := context.Background()
	//db, err := database.Connect(ctx, "postgresql://postgres:postgres@localhost:5432")
	db, err := database.Connect(ctx, os.Getenv("SUPABASE_DB"))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer db.Close()

	err = db.SetupTables(ctx)
	if err != nil {
		log.Fatalf("could not setup tables: %v", err)
	}

	token := os.Getenv("OPENAI_API_KEY")
	gpt := openai.NewClient(token)

	storage := storagego.NewClient(
		os.Getenv("SUPABASE_URL")+"/storage/v1",
		os.Getenv("SUPABASE_STORAGE_TOKEN"),
		nil)

	source := index.Index{
		DB:      db,
		GPT:     gpt,
		Storage: storage,
	}

	source.DB = db

	var messages []openai.ChatCompletionMessage

	// Read files recursively
	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if strings.Contains(path, ".git") {
			return nil
		}

		if strings.Contains(path, ".idea") {
			return nil
		}

		if !strings.Contains(path, "database") {
			return nil
		}

		if strings.HasSuffix(path, "test.go") {
			return nil
		}

		if !strings.HasSuffix(path, ".go") || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		byt, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("File: %v\nContent: %s", path, byt),
		})

		log.Printf("File: %v", path)

		return nil
	})

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: "Write tests for function CreateChat",
	})

	resp, err := gpt.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT4,
		Messages: messages,
		N:        1,
	})
	if err != nil {
		log.Fatal(err)
	}

	content := resp.Choices[0].Message.Content

	log.Printf("Usage: %+v", resp.Usage)

	err = os.WriteFile("database_test.go", []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
