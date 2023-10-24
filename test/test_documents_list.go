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

	ctx, userId := setup.createRandomSignInWithFunding(1000)
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

	setup.Report.Run("documents_list", func(t testing) bool {
		list, err := setup.documents.List(ctx, &pb.DocumentFilter{
			CollectionId: collection.Id,
		})
		if err != nil {
			return t.fail(err)
		}

		if len(list.Items) != 1 {
			return t.fail(fmt.Errorf("invalid document list: %v", list))
		}

		return t.pass()
	})

	setup.Report.Run("document_list_nothing", func(t testing) bool {
		list, err := setup.documents.List(ctx, &pb.DocumentFilter{
			Query:        "find nothing",
			CollectionId: collection.Id,
		})
		if err != nil {
			return t.fail(err)
		}

		if len(list.Items) != 0 {
			return t.fail(fmt.Errorf("invalid document list: %v", list))
		}

		return t.pass()
	})

	setup.Report.Run("document_list_query", func(t testing) bool {
		list, err := setup.documents.List(ctx, &pb.DocumentFilter{
			Query:        doc.Filename,
			CollectionId: collection.Id,
		})
		if err != nil {
			return t.fail(err)
		}

		if len(list.Items) != 1 {
			return t.fail(fmt.Errorf("invalid document list: %v", list))
		}

		if list.Items[0].Filename != doc.Filename {
			return t.fail(fmt.Errorf("invalid document name: %v", list.Items[0].Filename))
		}

		if list.Items[0].Id != doc.Id {
			return t.fail(fmt.Errorf("invalid document id: %v", list.Items[0].Id))
		}

		return t.pass()
	})

	setup.Report.Run("documents_list_wrong_collection", func(t testing) bool {
		list, err := setup.documents.List(ctx, &pb.DocumentFilter{
			Query:        doc.Filename,
			CollectionId: uuid.NewString(),
		})
		if err != nil {
			return t.fail(err)
		}

		if len(list.Items) != 0 {
			return t.fail(fmt.Errorf("invalid document list: %v", list))
		}

		return t.pass()
	})

	setup.Report.Run("documents_list_without_auth", func(t testing) bool {
		_, err = setup.documents.List(context.Background(), &pb.DocumentFilter{
			Query:        doc.Filename,
			CollectionId: collection.Id,
		})
		return t.expectError(err)
	})
}
