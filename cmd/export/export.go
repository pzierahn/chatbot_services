package main

import (
	"context"
	"encoding/json"
	"github.com/pzierahn/braingain/database"
	"log"
	"os"
)

type QdrantExport struct {
	Id        string
	Embedding []float32
	Filename  string
	Page      int
	Text      string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	old, err := database.Connect("localhost:6334")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() { _ = old.Close() }()

	results, err := old.GetAll(ctx, "DeSys")
	if err != nil {
		log.Fatal(err)
	}

	exports := make([]QdrantExport, len(results))

	log.Printf("results=%v", len(results))
	for inx, result := range results {
		text := result.Payload["content"].GetStringValue()
		filename := result.Payload["filename"].GetStringValue()
		page := result.Payload["page"].GetIntegerValue()

		// log.Printf("%v --> %d", filename, page)

		exports[inx] = QdrantExport{
			Id:        result.Id.GetUuid(),
			Embedding: result.Vectors.GetVector().Data,
			Filename:  filename,
			Page:      int(page),
			Text:      text,
		}
	}

	byt, err := json.MarshalIndent(exports, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("qdrant_export.json", byt, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
