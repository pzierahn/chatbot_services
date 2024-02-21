package test

import (
	"context"
	"fmt"
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

		meta := &pb.DocumentMetadata{
			Data: &pb.DocumentMetadata_Web{
				Web: &pb.Webpage{
					Title: "Wiki Penguin",
					Url:   "https://en.wikipedia.org/wiki/Penguin",
				},
			},
		}

		job := &pb.IndexJob{
			Id:           uuid.NewString(),
			CollectionId: collection.Id,
			Document:     meta,
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

		list, err := test.documents.List(ctx, &pb.DocumentFilter{
			CollectionId: collection.Id,
		})
		if err != nil {
			return err
		}

		web, ok := list.Items[job.Id]
		if !ok {
			return fmt.Errorf("document not found %v", list.Items)
		}

		if web.GetWeb().Title != meta.GetWeb().Title {
			return fmt.Errorf("title mismatch %v", web.GetWeb().Title)
		}

		if web.GetWeb().Url != meta.GetWeb().Url {
			return fmt.Errorf("url mismatch %v", web.GetWeb().Url)
		}

		return nil
	})
}
