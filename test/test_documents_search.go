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

	setup.Report.Run("documents_search_without_auth", func(t testing) bool {
		_, err = setup.documents.Search(context.Background(), &pb.SearchQuery{
			CollectionId: collection.Id,
		})
		return t.expectError(err)
	})

	setup.Report.Run("documents_search_without_query", func(t testing) bool {
		_, err = setup.documents.Search(ctx, &pb.SearchQuery{
			CollectionId: collection.Id,
		})
		return t.expectError(err)
	})

	setup.Report.Run("documents_search_limit", func(t testing) bool {
		results, err := setup.documents.Search(ctx, &pb.SearchQuery{
			Query:        "Clancy the Crab",
			CollectionId: collection.Id,
			Limit:        0,
		})
		if err != nil {
			return t.fail(err)
		}

		if len(results.Items) != 0 {
			return t.fail(fmt.Errorf("expected 0 results, got %d", len(results.Items)))
		}

		return t.pass()
	})

	setup.Report.Run("documents_search_valid", func(t testing) bool {
		results, err := setup.documents.Search(ctx, &pb.SearchQuery{
			Query:        "Clancy the Crab",
			CollectionId: collection.Id,
			Limit:        15,
		})
		if err != nil {
			return t.fail(err)
		}

		if len(results.Items) != 1 {
			return t.fail(fmt.Errorf("expected 1 results, got %d", len(results.Items)))
		}

		if results.Items[0].DocumentId != docId {
			return t.fail(fmt.Errorf("expected document id %s, got %s", docId, results.Items[0].DocumentId))
		}

		if results.Items[0].Score <= 0 {
			return t.fail(fmt.Errorf("expected score > 0, got %f", results.Items[0].Score))
		}

		return t.pass()
	})

	setup.Report.Run("documents_search_invalid", func(t testing) bool {
		results, err := setup.documents.Search(ctx, &pb.SearchQuery{
			Query:        "This is not in the document",
			CollectionId: collection.Id,
			Threshold:    0.8,
			Limit:        15,
		})
		if err != nil {
			return t.fail(err)
		}

		if len(results.Items) != 0 {
			return t.fail(fmt.Errorf("expected 0 results, got %d", len(results.Items)))
		}

		return t.pass()
	})
}
