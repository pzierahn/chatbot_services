package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
)

func (setup *Setup) ChatHistory() {

	ctx, userId := setup.createRandomSignIn()
	defer setup.DeleteUser(userId)

	collection, err := setup.collections.Create(ctx, &pb.Collection{Name: "test"})
	if err != nil {
		log.Fatal(err)
	}

	docId := uuid.NewString()
	path := fmt.Sprintf("%s/%s/%s.pdf", userId, collection.Id, docId)

	resp := setup.storage.UploadFile(bucket, path, bytes.NewReader(testPdf))
	if resp.Error != "" {
		log.Fatalf("upload failed: %v", resp.Error)
	}

	doc := &pb.Document{
		Id:           docId,
		CollectionId: collection.Id,
		Filename:     "test.pdf",
		Path:         path,
	}

	stream, err := setup.documents.Index(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, err = stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}

	prompt := &pb.Prompt{
		Prompt:       "What is the color of Clancy the Crab's house? Answer with one word",
		CollectionId: collection.Id,
		Options: &pb.PromptOptions{
			Model:       openai.GPT3Dot5Turbo,
			Temperature: 0,
			MaxTokens:   1,
			Threshold:   0,
			Limit:       1,
		},
	}

	setup.report.ExpectError("chat_history_without_auth", func() error {
		_, err = setup.chat.Chat(context.Background(), prompt)
		return err
	})

	setup.report.Run("chat_history_valid", func() error {
		response, err := setup.chat.Chat(ctx, prompt)
		if err != nil {
			return err
		}

		if response.Text != "red" {
			return fmt.Errorf("expected answer red, got %s", response.Text)
		}

		messages, err := setup.chat.GetChatMessages(ctx, collection)
		if err != nil {
			return err
		}

		if len(messages.Ids) != 1 {
			return fmt.Errorf("expected 1 message, got %d", len(messages.Ids))
		}

		if messages.Ids[0] != response.Id {
			return fmt.Errorf("expected message id %s, got %s", response.Id, messages.Ids[0])
		}

		return nil
	})

	setup.report.Run("chat_history_get_massage", func() error {
		response, err := setup.chat.Chat(ctx, prompt)
		if err != nil {
			return err
		}

		message, err := setup.chat.GetChatMessage(ctx, &pb.MessageID{Id: response.Id})
		if err != nil {
			return err
		}

		if message.Id != response.Id {
			return fmt.Errorf("expected message id %s, got %s", response.Id, message.Id)
		}

		if message.Text != response.Text {
			return fmt.Errorf("expected message text %s, got %s", response.Text, message.Text)
		}

		if message.CollectionId != response.CollectionId {
			return fmt.Errorf("expected message collection id %s, got %s", response.CollectionId, message.CollectionId)
		}

		if message.Prompt.Prompt != response.Prompt.Prompt {
			return fmt.Errorf("expected message prompt %v, got %v", response.Prompt.Prompt, message.Prompt.Prompt)
		}

		return nil
	})

	setup.report.Run("chat_history_wrong_collection", func() error {
		messages, err := setup.chat.GetChatMessages(ctx, &pb.Collection{
			Id: uuid.NewString(),
		})
		if err != nil {
			return err
		}

		if len(messages.Ids) != 0 {
			return fmt.Errorf("expected 0 messages, got %d", len(messages.Ids))
		}

		return nil
	})
}
