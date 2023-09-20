package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/database"
	"log"
)

type chatMessage struct {
	uid        string
	collection string
	prompt     string
	completion string
	pageIDs    []uuid.UUID
}

func (server *Server) storeChatMessage(ctx context.Context, message chatMessage) {

	dbMessage := database.ChatMessage{
		UID:        message.uid,
		Collection: message.collection,
		Prompt:     message.prompt,
		Completion: message.completion,
	}

	for _, doc := range message.pageIDs {
		dbMessage.Sources = append(dbMessage.Sources, database.ChatMessageSource{
			DocumentPage: doc,
		})
	}

	_, err := server.db.CreateChat(ctx, dbMessage)
	if err != nil {
		log.Printf("storeChatMessage: error %v", err)
		return
	}
}
