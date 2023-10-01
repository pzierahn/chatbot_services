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

	setup.report.ExpectError("documents_delete_invalid", func() error {
		_, err = setup.documents.Delete(ctx, &pb.Document{})
		return err
	})

	setup.report.ExpectError("documents_delete_without_auth", func() error {
		_, err = setup.documents.Delete(context.Background(), doc)
		return err
	})

	setup.report.Run("documents_delete", func() error {
		_, err = setup.documents.Delete(ctx, doc)
		if err != nil {
			return err
		}

		list, err := setup.documents.List(ctx, &pb.DocumentFilter{
			CollectionId: collection.Id,
		})
		if err != nil {
			return err
		}

		if len(list.Items) != 0 {
			return fmt.Errorf("invalid document list: %v", list)
		}

		files := setup.storage.ListFiles(bucket, userId+"/"+collection.Id, storage_go.FileSearchOptions{})
		if len(files) != 0 {
			return fmt.Errorf("invalid storage list: %v", files)
		}

		return nil
	})
}
