package test

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/pzierahn/brainboost/proto"
	"io"
	"log"
)

func (setup *Setup) DocumentsSearch() {

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

	setup.report.ExpectError("documents_search_without_auth", func() error {
		_, err = setup.documents.Search(context.Background(), &pb.SearchQuery{
			CollectionId: collection.Id,
		})
		return err
	})

	setup.report.ExpectError("documents_search_without_query", func() error {
		_, err = setup.documents.Search(ctx, &pb.SearchQuery{
			CollectionId: collection.Id,
		})
		return err
	})

	setup.report.Run("documents_search_limit", func() error {
		results, err := setup.documents.Search(ctx, &pb.SearchQuery{
			Query:        "Clancy the Crab",
			CollectionId: collection.Id,
			Limit:        0,
		})
		if err != nil {
			return err
		}

		if len(results.Items) != 0 {
			return fmt.Errorf("expected 0 results, got %d", len(results.Items))
		}

		return nil
	})

	setup.report.Run("documents_search_valid", func() error {
		results, err := setup.documents.Search(ctx, &pb.SearchQuery{
			Query:        "Clancy the Crab",
			CollectionId: collection.Id,
			Limit:        15,
		})
		if err != nil {
			return err
		}

		if len(results.Items) != 1 {
			return fmt.Errorf("expected 1 results, got %d", len(results.Items))
		}

		if results.Items[0].DocumentId != docId {
			return fmt.Errorf("expected document id %s, got %s", docId, results.Items[0].DocumentId)
		}

		if results.Items[0].Score <= 0 {
			return fmt.Errorf("expected score > 0, got %f", results.Items[0].Score)
		}

		return nil
	})

	setup.report.Run("documents_search_invalid", func() error {
		results, err := setup.documents.Search(ctx, &pb.SearchQuery{
			Query:        "This is not in the document",
			CollectionId: collection.Id,
			Limit:        15,
		})
		if err != nil {
			return err
		}

		if len(results.Items) != 0 {
			return fmt.Errorf("expected 0 results, got %d", len(results.Items))
		}

		return nil
	})
}
