package test

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/pzierahn/brainboost/proto"
	storage_go "github.com/supabase-community/storage-go"
	"io"
	"log"
)

func (setup *Setup) DocumentsDelete() {

	ctx, userId := setup.createRandomSignInWithFunding(1000)
	defer setup.DeleteUser(userId)

	collection, err := setup.collections.Create(ctx, &pb.Collection{Name: "test"})
	if err != nil {
		log.Fatal(err)
	}

	docId := uuid.NewString()
	path := fmt.Sprintf("%s/%s/%s.pdf", userId, collection.Id, docId)

	_, err = setup.storage.UploadFile(bucket, path, bytes.NewReader(testPdf))
	if err != nil {
		log.Fatalf("upload failed: %v", err)
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

	setup.Report.Run("documents_delete_invalid", func(t testing) bool {
		_, err = setup.documents.Delete(ctx, &pb.Document{})
		return t.expectError(err)
	})

	setup.Report.Run("documents_delete_without_auth", func(t testing) bool {
		_, err = setup.documents.Delete(context.Background(), doc)
		return t.expectError(err)
	})

	setup.Report.Run("documents_delete", func(t testing) bool {
		_, err = setup.documents.Delete(ctx, doc)
		if err != nil {
			return t.fail(err)
		}

		list, err := setup.documents.List(ctx, &pb.DocumentFilter{
			CollectionId: collection.Id,
		})
		if err != nil {
			return t.fail(err)
		}

		if len(list.Items) != 0 {
			return t.fail(fmt.Errorf("invalid document list: %v", list))
		}

		files, err := setup.storage.ListFiles(bucket, userId+"/"+collection.Id, storage_go.FileSearchOptions{})
		if err != nil {
			return t.fail(err)
		}

		if len(files) != 0 {
			return t.fail(fmt.Errorf("invalid storage list: %v", files))
		}

		return t.pass()
	})
}
