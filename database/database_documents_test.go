package database

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestUpsertDocument_Valid(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Create a test collection
	testCollection := setupTestCollection(t)

	// Insert the test collection into the database
	testCollectionID, err := testClient.CreateCollection(context.Background(), testCollection)
	if err != nil {
		t.Fatalf("Error upserting collection: %v", err)
	}

	testDocument := Document{
		UserId:     testCollection.UserId,
		Collection: testCollectionID,
		Filename:   "test.pdf",
		Path:       "/path/to/test.pdf",
		Pages: []*PageEmbedding{
			{Page: 1, Text: "Page 1 Content", Embedding: make([]float32, 1536)},
			{Page: 2, Text: "Page 2 Content", Embedding: make([]float32, 1536)},
		},
	}

	// Insert the test document into the database
	createdID, err := testClient.UpsertDocument(context.Background(), testDocument)
	if err != nil {
		t.Fatalf("Error upserting document: %v", err)
	}

	if createdID == uuid.Nil {
		t.Fatal("Expected createdID to be non-nil")
	}
}

func TestFindDocuments_Valid(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Create a test collection
	testCollection := setupTestCollection(t)

	testDocuments := []Document{
		{UserId: testCollection.UserId, Collection: testCollection.Id, Filename: "document1.pdf", Path: "/path/to/document1.pdf"},
		{UserId: testCollection.UserId, Collection: testCollection.Id, Filename: "document2.pdf", Path: "/path/to/document2.pdf"},
	}

	for _, doc := range testDocuments {
		_, err := testClient.UpsertDocument(context.Background(), doc)
		if err != nil {
			t.Fatalf("Error upserting document: %v", err)
		}
	}

	// Define a test query
	testQuery := DocumentQuery{
		UserId:     testCollection.UserId,
		Collection: testCollection.Id,
		Query:      "document",
	}

	// Find documents matching the query
	results, err := testClient.FindDocuments(context.Background(), testQuery)
	if err != nil {
		t.Fatalf("Error finding documents: %v", err)
	}

	// Ensure that the correct number of documents is retrieved
	if len(results) != len(testDocuments) {
		t.Fatalf("Expected %d documents, but got %d", len(testDocuments), len(results))
	}
}

func TestDeleteDocument_Valid(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Create a test collection
	testCollection := setupTestCollection(t)

	// Insert the test collection into the database
	testCollectionID, err := testClient.CreateCollection(context.Background(), testCollection)
	if err != nil {
		t.Fatalf("Error upserting collection: %v", err)
	}

	testDocument := Document{
		UserId:     testCollection.UserId,
		Collection: testCollectionID,
		Filename:   "document_to_delete.pdf",
		Path:       "/path/to/document_to_delete.pdf",
	}

	// Insert the test document into the database
	_, err = testClient.UpsertDocument(context.Background(), testDocument)
	if err != nil {
		t.Fatalf("Error upserting document: %v", err)
	}

	// Delete the test document
	err = testClient.DeleteDocument(context.Background(), testDocument.Id, testCollection.UserId)
	if err != nil {
		t.Fatalf("Error deleting document: %v", err)
	}
}

func TestUpdateDocumentName_Valid(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Create a test collection
	testCollection := setupTestCollection(t)
	testUserID := testCollection.UserId

	// Insert the test collection into the database
	testCollectionID, err := testClient.CreateCollection(context.Background(), testCollection)
	if err != nil {
		t.Fatalf("Error upserting collection: %v", err)
	}

	testDocument := Document{
		UserId:     testUserID,
		Collection: testCollectionID,
		Filename:   "old_filename.pdf",
		Path:       "/path/to/old_filename.pdf",
	}

	// Insert the test document into the database
	_, err = testClient.UpsertDocument(context.Background(), testDocument)
	if err != nil {
		t.Fatalf("Error upserting document: %v", err)
	}

	// Update the document name
	newFilename := "new_filename.pdf"
	testDocument.Filename = newFilename
	err = testClient.UpdateDocumentName(context.Background(), testDocument)
	if err != nil {
		t.Fatalf("Error updating document name: %v", err)
	}
}
