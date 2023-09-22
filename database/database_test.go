package database

import (
	"context"
	"testing"
)

// Initialize a test client and connection string
var testClient *Client

const testConnectionString = "postgres://patrick.zierahn:EMtDKkB0n4dP@ep-round-dream-25253463.eu-central-1.aws.neon.tech/neondb"

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
