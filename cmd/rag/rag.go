package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/braingain"
	"github.com/pzierahn/braingain/database"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"sort"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	ctx := context.Background()

	//conn, err := database.Connect("localhost:6334")
	conn, err := database.Connect(ctx, "postgresql://postgres:postgres@localhost:5432")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//search := "Explain the Decentralization, Scalability, and Consistency triangle."
	search := "Explain how DAG-Rider reaches consensus."

	token := os.Getenv("OPENAI_API_KEY")
	ai := openai.NewClient(token)

	chat := braingain.NewChat(conn, ai)
	chat.Model = openai.GPT4

	response, err := chat.RAG(ctx, search)
	if err != nil {
		log.Fatalf("Completion error: %v", err)
	}

	sources := make(map[uuid.UUID][]int)
	for _, source := range response.Documents {
		sources[source.Source] = append(sources[source.Source], source.Page)
	}

	keys := make([]uuid.UUID, 0, len(sources))
	for k := range sources {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Time() > keys[j].Time()
	})

	for _, id := range keys {
		pages := sources[id]

		doc, err := conn.GetDocument(ctx, id)
		if err != nil {
			log.Fatalf("GetDocument error: %v", err)
		}

		log.Printf("%s --> %v\n", doc.Filename, pages)
	}

	byt, _ := json.MarshalIndent(response.Costs, "", "  ")
	log.Printf("Costs: %s\n", string(byt))

	log.Println(response.Completion)
	_ = os.WriteFile("output.txt", []byte(response.Completion.Completion), 0644)

	byt, _ = json.MarshalIndent(response.Documents, "", "  ")
	_ = os.WriteFile("sources.json", byt, 0644)
}
