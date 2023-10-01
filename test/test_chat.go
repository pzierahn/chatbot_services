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
	"strings"
)

func (setup *Setup) ChatGenerate() {

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

	setup.report.ExpectError("chat_without_auth", func() error {
		_, err = setup.chat.Chat(context.Background(), prompt)
		return err
	})

	setup.report.Run("chat_auto_search", func() error {
		response, err := setup.chat.Chat(ctx, prompt)
		if err != nil {
			return err
		}

		if response.Text != "red" {
			return fmt.Errorf("expected answer red, got %s", response.Text)
		}

		return nil
	})

	setup.report.Run("chat_auto_with_pages", func() error {
		promptWithDoc := &pb.Prompt{
			Prompt:       "What is the color of Clancy the Crab's house? Answer with one word",
			CollectionId: collection.Id,
			Options: &pb.PromptOptions{
				Model:       openai.GPT3Dot5Turbo,
				Temperature: 0,
				MaxTokens:   10,
				Threshold:   0,
				Limit:       1,
			},
			Documents: []*pb.Prompt_Document{
				{
					Id:    doc.Id,
					Pages: []uint32{0},
				},
			},
		}

		response, err := setup.chat.Chat(ctx, promptWithDoc)
		if err != nil {
			return err
		}

		if !strings.Contains(response.Text, "red") {
			return fmt.Errorf("expected answer red, got %s", response.Text)
		}

		return nil
	})

	setup.report.Run("chat_auto_without_pages", func() error {
		promptWithDoc := &pb.Prompt{
			Prompt:       "What is the color of Clancy the Crab's house? Answer with one word",
			CollectionId: collection.Id,
			Options: &pb.PromptOptions{
				Model:       openai.GPT3Dot5Turbo,
				Temperature: 0,
				MaxTokens:   1,
				Threshold:   0,
				Limit:       1,
			},
			Documents: []*pb.Prompt_Document{
				{
					Id: doc.Id,
					// Add no pages to the document.
					Pages: []uint32{},
				},
			},
		}

		response, err := setup.chat.Chat(ctx, promptWithDoc)
		if err != nil {
			return err
		}

		if response.Text == "red" {
			return fmt.Errorf("didn't expect the right anwser")
		}

		if len(response.Documents[0].Pages) != 0 {
			return fmt.Errorf("expected page 0, got %d", response.Documents[0].Pages)
		}

		return nil
	})

	setup.report.Run("chat_auto_with_pages_invalid", func() error {
		promptWithDoc := &pb.Prompt{
			Prompt:       "What is the color of Clancy the Crab's house? Answer with one word",
			CollectionId: collection.Id,
			Options: &pb.PromptOptions{
				Model:       openai.GPT3Dot5Turbo,
				Temperature: 0,
				MaxTokens:   1,
				Threshold:   0,
				Limit:       1,
			},
			Documents: []*pb.Prompt_Document{
				{
					Id: doc.Id,
					Pages: []uint32{
						100,
					},
				},
			},
		}

		response, err := setup.chat.Chat(ctx, promptWithDoc)
		if err != nil {
			return err
		}

		if strings.ToLower(response.Text) == "red" {
			return fmt.Errorf("didn't expect the right anwser")
		}

		return nil
	})
}
