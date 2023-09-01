package main

import (
	"context"
	"encoding/json"
	"github.com/pzierahn/braingain/braingain"
	"github.com/pzierahn/braingain/database"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"sort"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	conn, err := database.Connect("localhost:6334")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() { _ = conn.Close() }()

	ctx := context.Background()

	search := "Differance between sharding and sidechains"

	token := os.Getenv("OPENAI_API_KEY")
	ai := openai.NewClient(token)

	chat := braingain.NewChat(conn, ai)
	//chat.Model = openai.GPT4

	response, err := chat.RAG(ctx, search)
	if err != nil {
		log.Fatalf("ChatCompletion error: %v", err)
	}

	sources := make(map[string][]int)
	for _, source := range response.Sources {
		sources[source.Filename] = append(sources[source.Filename], source.Page)
	}

	keys := make([]string, 0, len(sources))
	for k := range sources {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, filename := range keys {
		pages := sources[filename]
		log.Printf("%s --> %v\n", filename, pages)
	}

	byt, _ := json.MarshalIndent(response.Costs, "", "  ")
	log.Printf("Costs: %s\n", string(byt))

	log.Println(response.Completion)
	_ = os.WriteFile("output.txt", []byte(response.Completion), 0644)

	byt, _ = json.MarshalIndent(response.Sources, "", "  ")
	_ = os.WriteFile("sources.json", byt, 0644)
}
