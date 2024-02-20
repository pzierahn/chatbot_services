package test

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/pzierahn/chatbot_services/proto"
	"io"
)

func (test Tester) TestDocumentList() {
	test.runTest("TestDocumentList", func(ctx context.Context) error {
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

		list, err := test.documents.List(ctx, &pb.DocumentFilter{
			CollectionId: collection.Id,
		})
		if err != nil {
			return err
		}

		if len(list.Items) != 1 {
			return err
		}

		return nil
	})
}
