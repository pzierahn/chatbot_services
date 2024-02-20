package test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/pzierahn/chatbot_services/proto"
	"io"
)

func (test Tester) TestDocumentGet() {
	test.runTest("TestDocumentGet", func(ctx context.Context) error {
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

		// Get document
		doc, err := test.documents.Get(ctx, &pb.DocumentID{
			Id: job.Id,
		})
		if err != nil {
			return err
		}

		if doc.Id != job.Id {
			return fmt.Errorf("expected document id %s, got %s", job.Id, doc.Id)
		}

		if doc.CollectionId != collection.Id {
			return fmt.Errorf("expected collection id %s, got %s", collection.Id, doc.CollectionId)
		}

		switch doc.Metadata.Data.(type) {
		case *pb.DocumentMetadata_Web:
			web := doc.Metadata.GetWeb()
			if web.Url != "https://en.wikipedia.org/wiki/Penguin" {
				return fmt.Errorf("expected url %s, got %s", "https://en.wikipedia.org/wiki/Penguin", web.Url)
			}
			if web.Title != "Wiki Penguin" {
				return fmt.Errorf("expected title %s, got %s", "Wiki Penguin", web.Title)
			}
		default:
			return fmt.Errorf("expected web document, got %T", doc.Metadata.Data)
		}

		if len(doc.Chunks) == 0 {
			return fmt.Errorf("expected chunks, got none")
		}

		// Check if chunks are sorted
		for inx := 1; inx < len(doc.Chunks); inx++ {
			if doc.Chunks[inx].Index < doc.Chunks[inx-1].Index {
				return fmt.Errorf("chunks are not sorted")
			}
		}

		return nil
	})
}
