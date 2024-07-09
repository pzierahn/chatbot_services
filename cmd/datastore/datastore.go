package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/pdf"
	"log"
)

const (
	userId       = "j7jjxLD9rla2DrZoeUu3Tnft4812"
	threadId     = "bb05d2b7-47b7-4ea8-9a4e-b47ef5c99b79"
	collectionId = "173dd77e-681b-4f5e-a3b8-cb91f19a0f56"
)

//go:embed document.pdf
var document []byte

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()
	db, err := datastore.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new document
	pages, err := pdf.GetPagesFromBytes(ctx, document)
	if err != nil {
		log.Fatal(err)
	}

	chunks := make([]datastore.DocumentChunk, len(pages))
	for inx, page := range pages {
		chunks[inx] = datastore.DocumentChunk{
			Id:       uuid.New(),
			Text:     page,
			Position: inx,
		}
	}

	doc := &datastore.Document{
		Id:           uuid.New(),
		UserId:       userId,
		CollectionId: uuid.MustParse(collectionId),
		Name:         "document.pdf",
		Type:         datastore.DocumentTypePDF,
		Source:       "xxxx/document.pdf",
		Content:      chunks,
	}

	// Store the document
	err = db.InsertDocument(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the document
	doc, err = db.GetDocument(ctx, userId, doc.Id)
	if err != nil {
		log.Fatal(err)
	}

	byt, _ := json.MarshalIndent(doc, "", "  ")
	log.Printf("Document: %s", byt)
}
