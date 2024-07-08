package main

import (
	"context"
	"encoding/json"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/llm/voyageai"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := voyageai.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	embeddings, err := client.CreateEmbedding(ctx, &llm.EmbeddingRequest{
		Inputs: []string{"Hello World!"},
	})
	if err != nil {
		log.Fatal(err)
	}

	byt, _ := json.MarshalIndent(embeddings, "", "  ")
	log.Println(string(byt))
}
