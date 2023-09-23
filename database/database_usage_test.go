package database

import (
	"context"
	"github.com/google/uuid"
	"testing"
)

func TestCreateUsage(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Define a test Usage struct
	usage := Usage{
		UserID: uuid.New(),
		Model:  "TestModel",
		Input:  100,
		Output: 50,
	}

	// Call the CreateUsage function to insert a new usage record
	id, err := testClient.CreateUsage(context.Background(), usage)
	if err != nil {
		t.Errorf("CreateUsage returned an error: %v", err)
	}

	// Verify that the usage ID is not nil
	if id == uuid.Nil {
		t.Error("CreateUsage returned a nil UUID")
	}

	// Optionally, you can query the database to check if the record was inserted correctly.
}

func TestGetModelUsages(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Create some test usages
	usages := []Usage{
		{
			UserID: uuid.New(),
			Model:  "TestModel",
			Input:  100,
			Output: 50,
		},
		{
			UserID: uuid.New(),
			Model:  "TestModel",
			Input:  200,
			Output: 100,
		},
	}

	// Insert the test usages into the database
	for _, usage := range usages {
		_, err := testClient.CreateUsage(context.Background(), usage)
		if err != nil {
			t.Errorf("CreateUsage returned an error: %v", err)
		}
	}

	for _, usage := range usages {
		// Call the GetModelUsages function to retrieve the usages
		retrievedUsages, err := testClient.GetModelUsages(context.Background(), usage.UserID)
		if err != nil {
			t.Errorf("GetModelUsages returned an error: %v", err)
		}

		// Verify that the usages were retrieved correctly
		if retrievedUsages[0].Model != usage.Model {
			t.Errorf("GetModelUsages returned an incorrect Model: %v", retrievedUsages[0].Model)
		}

		if retrievedUsages[0].Input != usage.Input {
			t.Errorf("GetModelUsages returned an incorrect Input: %v", retrievedUsages[0].Input)
		}

		if retrievedUsages[0].Output != usage.Output {
			t.Errorf("GetModelUsages returned an incorrect Output: %v", retrievedUsages[0].Output)
		}
	}

}
