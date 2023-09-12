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
	"strings"
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

	path := baseDir + "/Lecture Slides/"
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".pdf") {
			continue
		}

		log.Printf("Processing: %v", file.Name())

		byt, err := os.ReadFile(path + file.Name())
		if err != nil {
			log.Fatalf("could not read file: %v", err)
		}

		doc := index.DocumentId{
			UserId:     uuid.MustParse("50372462-3137-4ed9-9950-ad033fa24bfc"),
			Collection: uuid.MustParse("b452f76d-c1e4-4cdb-979f-08a4521d3372"),
			Filename:   file.Name(),
		}

		id, err := source.Process(ctx, doc, byt)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Success! %v", id)
	}

	//file := baseDir + "/Further Readings/IPTPS2002.pdf"
	//byt, err := os.ReadFile(file)
	//if err != nil {
	//	log.Fatalf("could not read file: %v", err)
	//}
	//
	//doc := index.DocumentId{
	//	UserId:     uuid.MustParse("3bc23192-230a-4366-b8ec-0bd7cce69510"),
	//	Collection: uuid.MustParse("b452f76d-c1e4-4cdb-979f-08a4521d3372"),
	//	Filename:   filepath.Base(file),
	//}
	//
	//id, err := source.Process(ctx, doc, byt)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Printf("Success! %v", id)
}
