package test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/pzierahn/chatbot_services/proto"
	"io"
)

func (test Tester) TestDocumentDelete() {
	test.runTest("TestDocumentDelete", func(ctx context.Context) error {
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
						Url: "https://en.wikipedia.org/wiki/Penguin",
					},
				},
			},
		}

		stream, err := test.documents.IndexDocument(ctx, job)
		if err != nil {
			return err
		}

		for {
			_, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
		}

		_, err = test.documents.Delete(ctx, &pb.DocumentID{
			Id: job.Id,
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

		if _, ok := list.Items[job.Id]; ok {
			return fmt.Errorf("document not deleted")
		}

		return nil
	})
}
