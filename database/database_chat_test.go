package database

import (
	"context"
	"github.com/google/uuid"
	"sort"
	"testing"
	"time"
)

func TestGetChatMessages_Valid(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Create a test collection
	testCollection := setupTestCollection(t)
	testUserID := testCollection.UserID

	// Create history chat messages
	testChatMessages := []ChatMessage{
		{
			UserID:       testUserID,
			CollectionID: testCollection.ID,
			Prompt:       "Test Prompt 1",
			Completion:   "Test Completion 1",
		},
		{
			UserID:       testUserID,
			CollectionID: testCollection.ID,
			Prompt:       "Test Prompt 2",
			Completion:   "Test Completion 2",
		},
	}

	// Store the created IDs of chat messages
	var createdChatMessageIDs []uuid.UUID

	// Insert chat messages into the database
	for _, message := range testChatMessages {
		id, err := testClient.CreateChat(context.Background(), message)
		if err != nil {
			t.Fatalf("Failed to create chat message: %v", err)
		}
		createdChatMessageIDs = append(createdChatMessageIDs, id)
	}

	// Call GetChatMessages function
	ids, err := testClient.GetChatMessages(context.Background(), testUserID.String(), testCollection.ID.String())
	if err != nil {
		t.Fatalf("Failed to get chat messages: %v", err)
	}

	// Assert that we retrieve the correct number of messages
	if len(ids) != len(testChatMessages) {
		t.Fatalf("Expected %d chat messages, got %d", len(testChatMessages), len(ids))
	}

	sort.Slice(ids, func(i, j int) bool {
		return ids[i].String() < ids[j].String()
	})
	sort.Slice(createdChatMessageIDs, func(i, j int) bool {
		return createdChatMessageIDs[i].String() < createdChatMessageIDs[j].String()
	})

	// Assert that IDs are the same as the ones we created
	for i := range ids {
		if ids[i] != createdChatMessageIDs[i] {
			t.Fatalf("Expected ID %s, got ID %s", createdChatMessageIDs[i], ids[i])
		}
	}
}

func TestGetChatMessages_NoChats(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Call GetChatMessages function for user with no chat messages
	ids, err := testClient.GetChatMessages(context.Background(), uuid.New().String(), uuid.New().String())
	if err != nil {
		t.Fatalf("Failed to get non-existent chat messages: %v", err)
	}

	// Assert that no chat messages are retrieved
	if len(ids) != 0 {
		t.Fatalf("Expected to get 0 chat message, but got %d", len(ids))
	}
}

func TestGetChatMessage_Valid(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Create a test collection
	testCollection := setupTestCollection(t)

	testMessage := ChatMessage{
		UserID:       testCollection.UserID,
		CollectionID: testCollection.ID,
		Prompt:       "Test Prompt",
		Completion:   "Test Completion",
		CreateAt:     time.Now(),
	}

	// Insert the test message into the database
	id, err := testClient.CreateChat(context.Background(), testMessage)
	if err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	// Retrieve the message
	retrievedMessage, err := testClient.GetChatMessage(context.Background(), id.String(), testCollection.UserID.String())
	if err != nil {
		t.Fatalf("Failed to retrieve message: %v", err)
	}

	// Assert the contents of the retrieved message
	if retrievedMessage.ID != id {
		t.Errorf("Expected ID to be %s, but was %s", id, retrievedMessage.ID)
	}

	if retrievedMessage.UserID != testMessage.UserID {
		t.Errorf("Expected UserID to be %s, but was %s", testMessage.UserID, retrievedMessage.UserID)
	}

	if retrievedMessage.CollectionID != testMessage.CollectionID {
		t.Errorf("Expected CollectionID to be %s, but was %s", testMessage.CollectionID, retrievedMessage.CollectionID)
	}

	if retrievedMessage.Prompt != testMessage.Prompt {
		t.Errorf("Expected Prompt to be %s, but was %s", testMessage.Prompt, retrievedMessage.Prompt)
	}

	if retrievedMessage.Completion != testMessage.Completion {
		t.Errorf("Expected Completion to be %s, but was %s", testMessage.Completion, retrievedMessage.Completion)
	}

}

func TestGetChatMessage_Invalid(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Attempt to retrieve a non-existent message
	_, err := testClient.GetChatMessage(context.Background(), uuid.New().String(), uuid.New().String())
	if err == nil {
		t.Fatalf("Expected error when retrieving non-existent message, but got nil")
	}
}

func TestCreateChat_Valid(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	testCollection := setupTestCollection(t)
	//testDoc := setupTestDocument(t, testCollection)

	// Define a test chat message
	testChatMessage := ChatMessage{
		UserID:       testCollection.UserID,
		CollectionID: testCollection.ID,
		Prompt:       "Test Prompt",
		Completion:   "Test Completion",
		Sources: []ChatMessageSource{
			{
				DocumentPage: uuid.New(),
				Filename:     "Test Filename",
				Page:         1,
			},
		},
	}

	// Call the function under test
	id, err := testClient.CreateChat(context.Background(), testChatMessage)
	if err != nil {
		t.Errorf("CreateChat returned an error: %v", err)
	}

	// Verify that the chat message ID is not nil
	if id == uuid.Nil {
		t.Error("CreateChat returned a nil UUID")
	}

	// Optionally, you can query the database to check if the record was inserted correctly.
}

func TestCreateChat_Invalid(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Define an invalid chat message
	invalidChatMessage := ChatMessage{
		// The UserID field is intentionally left empty to make this chat message invalid
		CollectionID: uuid.New(),
		Prompt:       "Test Prompt",
		Completion:   "Test Completion",
	}

	// Call the function under test
	id, err := testClient.CreateChat(context.Background(), invalidChatMessage)
	if err == nil {
		t.Error("Expected CreateChat to return an error, but it did not")
	}

	// Verify that the chat message ID is nil
	if id != uuid.Nil {
		t.Errorf("Expected CreateChat to return a nil UUID, but it did not")
	}
}
