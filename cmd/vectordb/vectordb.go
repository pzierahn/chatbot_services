package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/llm/openai"
	"github.com/pzierahn/chatbot_services/search"
	pineconesearch "github.com/pzierahn/chatbot_services/search/pinecone"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	openaiClient, err := openai.New()
	if err != nil {
		log.Fatalf("failed to create openai client: %v", err)
	}

	pc, err := pineconesearch.New(openaiClient, "chatbot")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully created a new Pinecone search object!")

	ctx := context.Background()

	fragment := &search.Fragment{
		Id:           uuid.NewString(),
		Text:         "Ich bin ein Test",
		UserId:       uuid.NewString(),
		DocumentId:   uuid.NewString(),
		CollectionId: uuid.NewString(),
		Position:     0,
	}

	usage, err := pc.Upsert(ctx, []*search.Fragment{fragment})
	if err != nil {
		log.Fatalf("failed to upsert fragment: %v", err)
	}

	log.Printf("Successfully upserted fragment with usage: %v", usage)
}
