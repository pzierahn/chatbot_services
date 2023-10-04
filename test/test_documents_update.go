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

func (setup *Setup) DocumentsUpdate() {

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

	setup.Report.Run("documents_update_without_auth", func(t testing) bool {
		_, err = setup.documents.Update(context.Background(), doc)
		return t.expectError(err)
	})

	setup.Report.Run("documents_update_without_doc", func(t testing) bool {
		_, err = setup.documents.Update(ctx, &pb.Document{})
		return t.expectError(err)
	})

	setup.Report.Run("documents_update_valid", func(t testing) bool {
		update := &pb.Document{
			Id:           docId,
			CollectionId: collection.Id,
			Filename:     "Updated-Filename.pdf",
		}

		_, err = setup.documents.Update(ctx, update)
		if err != nil {
			return t.fail(err)
		}

		list, err := setup.documents.List(ctx, &pb.DocumentFilter{
			Query:        update.Filename,
			CollectionId: collection.Id,
		})
		if err != nil {
			return t.fail(err)
		}

		if len(list.Items) != 1 {
			return t.fail(fmt.Errorf("expected 1 document, got %d", len(list.Items)))
		}

		if list.Items[0].Filename != update.Filename {
			return t.fail(fmt.Errorf("expected filename %s, got %s", update.Filename, list.Items[0].Filename))
		}

		if list.Items[0].Id != doc.Id {
			return t.fail(fmt.Errorf("expected id %s, got %s", doc.Id, list.Items[0].Id))
		}

		return t.pass()
	})
}
