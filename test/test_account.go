package test

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
)

func (setup *Setup) Account() {

	ctx, userId := setup.createRandomSignInWithFunding(1000)
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

	setup.Report.Run("account_without_auth", func(t testing) bool {
		_, err = setup.account.GetModelUsages(context.Background(), &emptypb.Empty{})
		return t.expectError(err)
	})

	setup.Report.Run("account_usage", func(t testing) bool {
		_, err = setup.chat.Chat(ctx, &pb.Prompt{
			Prompt:       "Write a test sentence.",
			CollectionId: collection.Id,
			Options: &pb.PromptOptions{
				Model:     openai.GPT3Dot5Turbo,
				MaxTokens: 10,
			},
		})
		if err != nil {
			return t.fail(err)
		}

		usage, err := setup.account.GetModelUsages(ctx, &emptypb.Empty{})
		if err != nil {
			return t.fail(err)
		}

		if len(usage.Items) == 0 {
			return t.fail(fmt.Errorf("no usage generated"))
		}

		for _, item := range usage.Items {
			if item.Model == "" {
				return t.fail(fmt.Errorf("invalid model name: %v", item.Model))
			}

			if item.Input <= 0 {
				return t.fail(fmt.Errorf("invalid input tokens: %v", item.Input))
			}

			if item.Model == openai.GPT3Dot5Turbo && item.Output <= 0 {
				return t.fail(fmt.Errorf("invalid output tokens: %v", item.Output))
			}
		}

		return t.pass()
	})

}
