package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/database"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

type chatMessage struct {
	uid        uuid.UUID
	collection uuid.UUID
	prompt     string
	completion string
	pageIDs    []uuid.UUID
}

func (server *Server) storeChatMessage(ctx context.Context, message chatMessage) {

	dbMessage := database.ChatMessage{
		UserID:       message.uid,
		CollectionID: message.collection,
		Prompt:       message.prompt,
		Completion:   message.completion,
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

func (server *Server) GetChatMessages(ctx context.Context, collection *pb.Collection) (*pb.ChatMessages, error) {
	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	ids, err := server.db.GetChatMessages(ctx, uid.String(), collection.Id)
	if err != nil {
		return nil, err
	}

	messages := &pb.ChatMessages{}
	for _, id := range ids {
		messages.Ids = append(messages.Ids, id.String())
	}

	return messages, nil
}

func (server *Server) GetChatMessage(ctx context.Context, id *pb.MessageID) (*pb.ChatMessage, error) {
	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	chatId, err := uuid.Parse(id.Id)
	if err != nil {
		return nil, err
	}

	message, err := server.db.GetChatMessage(ctx, chatId, uid.String())
	if err != nil {
		return nil, err
	}

	completion := &pb.ChatMessage{
		Prompt: &pb.Prompt{
			Prompt: message.Prompt,
		},
		Text:      message.Completion,
		Timestamp: timestamppb.New(*message.CreateAt),
	}

	filename := make(map[uuid.UUID]string)
	pages := make(map[uuid.UUID][]uint32)

	for _, source := range message.Sources {
		filename[source.ID] = source.Filename
		pages[source.ID] = append(pages[source.ID], uint32(source.Page))
	}

	for key := range filename {
		completion.Documents = append(completion.Documents, &pb.ChatMessage_Document{
			Id:       key.String(),
			Filename: filename[key],
			Pages:    pages[key],
		})
	}

	return completion, nil
}
