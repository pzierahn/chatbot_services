package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/llm"
	"log"
	"time"
)

const (
	userId   = "j7jjxLD9rla2DrZoeUu3Tnft4812"
	threadId = "bb05d2b7-47b7-4ea8-9a4e-b47ef5c99b79"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()
	db, err := datastore.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//
	// Insert messages into the database.
	//

	callId := uuid.NewString()
	err = db.StoreThread(ctx, &datastore.Thread{
		Id:           uuid.New(),
		ThreadId:     uuid.MustParse(threadId),
		UserId:       userId,
		CollectionId: uuid.New(),
		Timestamp:    time.Now(),
		Messages: []*llm.Message{{
			Role:    llm.RoleUser,
			Content: "Hallo Bot!",
		}, {
			Role: llm.RoleAssistant,
			ToolCalls: []llm.ToolCall{{
				CallID: callId,
				Function: llm.Function{
					Name:      "echo",
					Arguments: "{\"text\": \"Hallo User!\"}",
				},
			}},
		}, {
			Role: llm.RoleUser,
			ToolResponses: []llm.ToolResponse{{
				CallID:  callId,
				Content: "{\"text\": \"Hallo User!\"}",
			}},
		},
		}})
	if err != nil {
		log.Fatal(err)
	}

	//
	// Get messages from the database.
	//

	results, err := db.GetThread(ctx, userId, uuid.MustParse(threadId))
	if err != nil {
		log.Fatal(err)
	}

	byt, _ := json.MarshalIndent(results, "", "  ")
	log.Println("thread:", string(byt))
}
