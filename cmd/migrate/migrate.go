package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/database"
	"log"
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

	pgv, err := database.Connect(ctx, "postgresql://postgres:postgres@localhost:5432")
	if err != nil {
		log.Fatal(err)
	}

	err = pgv.CreateExtension(ctx)
	if err != nil {
		log.Fatal(err)
	}

	_ = pgv.DropTables(ctx)

	err = pgv.CreateTables(ctx)
	if err != nil {
		log.Fatal(err)
	}

	sourceId := make(map[string]uuid.UUID)
	for _, export := range exports {
		if _, ok := sourceId[export.Filename]; ok {
			continue
		}

		if export.Filename == "" {
			continue
		}

		id, err := pgv.CreateSource(ctx, database.Document{
			Filename: export.Filename,
		})
		if err != nil {
			log.Fatal(err)
		}

		sourceId[export.Filename] = id
		log.Printf("%v --> %v", export.Filename, id)
	}

	for _, export := range exports {
		if export.Filename == "" {
			continue
		}

		id, ok := sourceId[export.Filename]
		if !ok {
			log.Fatalf("source not found: %v", export.Filename)
		}

		docId := uuid.MustParse(export.Id)

		_, err := pgv.Upsert(ctx, database.Point{
			Id:        &docId,
			Source:    id,
			Embedding: export.Embedding,
			Page:      export.Page,
			Text:      export.Text,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
