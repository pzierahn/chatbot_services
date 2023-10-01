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

func (setup *Setup) DocumentsList() {

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

	setup.report.Run("documents_list", func() error {
		list, err := setup.documents.List(ctx, &pb.DocumentFilter{
			CollectionId: collection.Id,
		})
		if err != nil {
			return err
		}

		if len(list.Items) != 1 {
			return fmt.Errorf("invalid document list: %v", list)
		}

		return nil
	})

	setup.report.Run("document_list_nothing", func() error {
		list, err := setup.documents.List(ctx, &pb.DocumentFilter{
			Query:        "find nothing",
			CollectionId: collection.Id,
		})
		if err != nil {
			return err
		}

		if len(list.Items) != 0 {
			return fmt.Errorf("invalid document list: %v", list)
		}

		return nil
	})

	setup.report.Run("document_list_query", func() error {
		list, err := setup.documents.List(ctx, &pb.DocumentFilter{
			Query:        doc.Filename,
			CollectionId: collection.Id,
		})
		if err != nil {
			return nil
		}

		if len(list.Items) != 1 {
			return fmt.Errorf("invalid document list: %v", list)
		}

		return nil
	})

	setup.report.Run("documents_list_wrong_collection", func() error {
		list, err := setup.documents.List(ctx, &pb.DocumentFilter{
			Query:        doc.Filename,
			CollectionId: uuid.NewString(),
		})
		if err != nil {
			return err
		}

		if len(list.Items) != 0 {
			return fmt.Errorf("invalid document list: %v", list)
		}

		return err
	})

	setup.report.ExpectError("documents_list_without_auth", func() error {
		_, err = setup.documents.List(context.Background(), &pb.DocumentFilter{
			Query:        doc.Filename,
			CollectionId: collection.Id,
		})
		return err
	})
}
