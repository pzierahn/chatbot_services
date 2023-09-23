package database

import (
	"context"
	"os"
	"testing"
)

// Initialize a test client and connection string
var testClient *Client

var testConnectionString = os.Getenv("TEST_DATABASE")

func setupTestClient(t *testing.T) {
	if testClient != nil {
		t.Fatal("Test client already initialized")
	}
	client, err := Connect(context.Background(), testConnectionString)
	if err != nil {
		t.Fatalf("Error connecting to the database: %v", err)
	}
	err = client.SetupTables(context.Background())
	if err != nil {
		t.Fatalf("Error setup tables: %v", err)
	}

	testClient = client
}

func tearDownTestClient(t *testing.T) {
	if testClient != nil {
		err := testClient.DropTables(context.Background())
		if err != nil {
			t.Fatalf("Error dropping tables: %v", err)
		}

		testClient.Close()
		testClient = nil
	}
}
