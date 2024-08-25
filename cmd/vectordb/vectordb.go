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
		UserId:       "d643dc64-e3f3-4662-b2a6-cdd2769ddcfc",
		DocumentId:   "374cb01a-7c58-4665-a5b9-823781915c04",
		CollectionId: "38457bd8-1628-4248-a95e-19718232c078",
		Text:         "Ich bin ein Test 0",
		Position:     0,
	}

	usage, err := pc.Upsert(ctx, []*search.Fragment{fragment})
	if err != nil {
		log.Fatalf("failed to upsert fragment: %v", err)
	}

	log.Printf("Successfully upserted fragment with usage: %v", usage)

	//results, err := pc.Search(ctx, search.Query{
	//	UserId:       "a643baf0-2bca-42dc-b254-251cb2fcd78e",
	//	CollectionId: "560424b3-50ac-490f-8b80-e506028f7d2f",
	//	Query:        "Test",
	//	Limit:        10,
	//	Threshold:    0.3,
	//})
	//if err != nil {
	//	log.Fatalf("failed to search: %v", err)
	//}
	//
	//log.Printf("Successfully searched for fragments: %v", results)

	//err = pc.DeleteDocument(ctx,
	//	"d643dc64-e3f3-4662-b2a6-cdd2769ddcfc",
	//	"38457bd8-1628-4248-a95e-19718232c079",
	//	"374cb01a-7c58-4665-a5b9-823781915c05")
	//if err != nil {
	//	log.Fatalf("failed to delete document: %v", err)
	//}

	//err = pc.DeleteCollection(ctx, "d643dc64-e3f3-4662-b2a6-cdd2769ddcfc", "38457bd8-1628-4248-a95e-19718232c079")
	//if err != nil {
	//	log.Fatalf("failed to delete document: %v", err)
	//}
}
