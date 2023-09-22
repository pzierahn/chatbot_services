package database

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func setupTestCollection(t *testing.T) *Collection {
	testUserID := uuid.New()
	testCollection := &Collection{
		UserId: testUserID,
		Name:   "Test Collection",
	}

	createdID, err := testClient.CreateCollection(context.Background(), testCollection)
	if err != nil {
		t.Fatalf("Error creating collection: %v", err)
	}

	if createdID == uuid.Nil {
		t.Fatal("Expected createdID to be non-nil")
	}

	return testCollection
}

func TestCreateCollection_Valid(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	testCollection := setupTestCollection(t)

	createdID, err := testClient.CreateCollection(context.Background(), testCollection)
	if err != nil {
		t.Fatalf("Error creating collection: %v", err)
	}

	if createdID == uuid.Nil {
		t.Fatal("Expected createdID to be non-nil")
	}
}

func TestListCollections_Valid(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Insert test collections into the database
	testUserID := uuid.New()
	testCollections := []*Collection{
		{UserId: testUserID, Name: "Collection 1"},
		{UserId: testUserID, Name: "Collection 2"},
		{UserId: uuid.New(), Name: "Collection 3"}, // Different user ID
	}

	for _, coll := range testCollections {
		_, err := testClient.CreateCollection(context.Background(), coll)
		if err != nil {
			t.Fatalf("Error creating collection: %v", err)
		}
		defer func(testClient *Client, ctx context.Context, coll *Collection) {
			err := testClient.DeleteCollection(ctx, coll)
			if err != nil {
				t.Fatalf("Error deleting collection: %v", err)
			}
		}(testClient, context.Background(), coll)
	}

	// Retrieve collections for the test user
	collections, err := testClient.ListCollections(context.Background(), testUserID)
	if err != nil {
		t.Fatalf("Error listing collections: %v", err)
	}

	// Ensure that the correct number of collections is retrieved
	if len(collections) != 2 {
		t.Fatalf("Expected 2 collections, but got %d", len(collections))
	}
}

func TestListCollections_InvalidUser(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Attempt to list collections for a non-existent user
	nonExistentUserID := uuid.New()

	collections, err := testClient.ListCollections(context.Background(), nonExistentUserID)
	if err != nil {
		t.Fatalf("Error listing collections: %v", err)
	}

	// Ensure that no collections are retrieved for the non-existent user
	if len(collections) != 0 {
		t.Fatalf("Expected 0 collections, but got %d", len(collections))
	}
}
