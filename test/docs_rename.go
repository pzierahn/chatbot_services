package test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/pzierahn/chatbot_services/proto"
	"io"
)

func (test Tester) TestDocumentRename() {
	test.runTest("TestDocumentRename", func(ctx context.Context) error {
		collection, err := test.collections.Create(ctx, &pb.Collection{
			Name: "test",
		})
		if err != nil {
			return err
		}

		job := &pb.IndexJob{
			Id:           uuid.NewString(),
			CollectionId: collection.Id,
			Document: &pb.DocumentMetadata{
				Data: &pb.DocumentMetadata_Web{
					Web: &pb.Webpage{
						Title: "Wiki Penguin",
						Url:   "https://en.wikipedia.org/wiki/Penguin",
					},
				},
			},
		}

		stream, err := test.documents.Index(ctx, job)
		if err != nil {
			return err
		}

		for {
			_, err = stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
		}

		const newTitle = "New title"
		_, err = test.documents.Rename(ctx, &pb.RenameDocument{
			Id: job.Id,
			RenameTo: &pb.RenameDocument_WebpageTitle{
				WebpageTitle: newTitle,
			},
		})
		if err != nil {
			return err
		}

		list, err := test.documents.List(ctx, &pb.DocumentFilter{
			CollectionId: collection.Id,
		})
		if err != nil {
			return err
		}

		doc, ok := list.Items[job.Id]
		if !ok {
			return fmt.Errorf("document not deleted")
		}

		if doc.GetWeb().Title != newTitle {
			return fmt.Errorf("title not changed")
		}

		return nil
	})
}
