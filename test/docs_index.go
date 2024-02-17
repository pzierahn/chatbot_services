package test

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"io"
	"log"
)

func (test Tester) TestWebpageIndex() {
	test.runTest("TestWebpageIndex", func(ctx context.Context) error {
		collection, err := test.collections.Create(ctx, &pb.Collection{
			Name: "test",
		})
		if err != nil {
			return err
		}

		doc := &pb.IndexJob{
			CollectionId: collection.Id,
			Document: &pb.DocumentMetadata{
				Data: &pb.DocumentMetadata_Web{
					Web: &pb.Webpage{
						Url: "https://en.wikipedia.org/wiki/Penguin",
						//Url: "https://fuks.org",
					},
				},
			},
		}

		stream, err := test.documents.IndexDocument(ctx, doc)
		if err != nil {
			return err
		}

		for {
			status, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			log.Printf("status: %v", status.Status)
		}

		return nil
	})
}
