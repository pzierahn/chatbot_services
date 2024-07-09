package main

import (
	"context"
	"encoding/json"
	"github.com/pzierahn/chatbot_services/search"
	"github.com/pzierahn/chatbot_services/search/qdrant"
	"log"
	"os"
)

const (
	collectionId = "5feaef59-2dff-430c-9f44-3b5d24f25b54"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := os.Setenv("CHATBOT_QDRANT_INSECURE", "true")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv("CHATBOT_QDRANT_URL", "localhost:6334")
	if err != nil {
		log.Fatal(err)
	}

	db, err := qdrant.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	//err = db.Upsert(ctx, []*vectordb.Fragment{{
	//	Id:           uuid.NewString(),
	//	DocumentId:   uuid.NewString(),
	//	UserId:       uuid.NewString(),
	//	CollectionId: collectionId,
	//	Text:         "Ich habe die Haare schön.",
	//}})
	//if err != nil {
	//	log.Fatal(err)
	//}

	results, err := db.Search(ctx, search.Query{
		CollectionId: collectionId,
		Query:        "Haare",
		Limit:        10,
		Threshold:    0.25,
	})
	if err != nil {
		log.Fatal(err)
	}

	byt, _ := json.MarshalIndent(results, "", "  ")
	log.Println(string(byt))
}
