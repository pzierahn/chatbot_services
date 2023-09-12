package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/database"
	"github.com/pzierahn/braingain/index"
	"github.com/sashabaranov/go-openai"
	storagego "github.com/supabase-community/storage-go"
	"log"
	"os"
	"path/filepath"
)

const (
	baseDir = "/Users/patrick/patrick.zierahn@gmail.com - Google Drive/My Drive/KIT/2023-SS/DeSys/"
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

	token := os.Getenv("OPENAI_API_KEY")
	gpt := openai.NewClient(token)

	storage := storagego.NewClient(
		"https://fikepkklraklkitnlfxi.supabase.co/storage/v1",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImZpa2Vwa2tscmFrbGtpdG5sZnhpIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTY5MzA0MjMxNSwiZXhwIjoyMDA4NjE4MzE1fQ.dfPSnBx4dKIeTcT4XCj5moufYGQXWajjfzeDct4tLSA",
		nil)

	source := index.Index{
		DB:      db,
		GPT:     gpt,
		Storage: storage,
	}

	//file := baseDir + "/Further Readings/IPTPS2002.pdf"
	file := baseDir + "/Lecture Slides/DeSys_04_Leader_Election.pdf"
	byt, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("could not read file: %v", err)
	}

	doc := index.DocumentId{
		UserId:     uuid.MustParse("50372462-3137-4ed9-9950-ad033fa24bfc"),
		Collection: uuid.MustParse("b452f76d-c1e4-4cdb-979f-08a4521d3372"),
		Filename:   filepath.Base(file),
	}

	id, err := source.Process(ctx, doc, byt)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Success! %v", id)
}
