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

	setup.Report.Run("chat_without_auth", func(t testing) bool {
		_, err = setup.chat.Chat(context.Background(), prompt)
		return t.expectError(err)
	})

	setup.Report.Run("chat_auto_search", func(t testing) bool {
		response, err := setup.chat.Chat(ctx, prompt)
		if err != nil {
			return t.fail(err)
		}

		if response.Text != "red" {
			return t.fail(fmt.Errorf("expected answer red, got %s", response.Text))
		}

		return t.pass()
	})

	setup.Report.Run("chat_auto_with_pages", func(t testing) bool {
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
			return t.fail(err)
		}

		if !strings.Contains(response.Text, "red") {
			return t.fail(fmt.Errorf("expected answer red, got %s", response.Text))
		}

		return t.pass()
	})

	setup.Report.Run("chat_auto_without_pages", func(t testing) bool {
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
			return t.fail(err)
		}

		if response.Text == "red" {
			return t.fail(fmt.Errorf("didn't expect the right anwser"))
		}

		if len(response.Documents[0].Pages) != 0 {
			return t.fail(fmt.Errorf("expected page 0, got %d", response.Documents[0].Pages))
		}

		return t.pass()
	})

	setup.Report.Run("chat_auto_with_pages_invalid", func(t testing) bool {
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
			return t.fail(err)
		}

		if strings.ToLower(response.Text) == "red" {
			return t.fail(fmt.Errorf("didn't expect the right anwser"))
		}

		return t.pass()
	})
}
