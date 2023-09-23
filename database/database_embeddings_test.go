package database

import (
	"context"
	"github.com/google/uuid"
	"testing"
	// Import any necessary testing libraries or mocks
)

func TestSearch(t *testing.T) {
	setupTestClient(t)
	defer tearDownTestClient(t)

	// Create a test collection
	testCollection := setupTestCollection(t)

	onesEmbedding := make([]float32, 1536)
	otherEmbedding := make([]float32, 1536)
	for i := range onesEmbedding {
		onesEmbedding[i] = 1.0
		otherEmbedding[i] = -1.0
	}

	// Create a SearchQuery for testing
	query := SearchQuery{
		UserID:       testCollection.UserID, // Replace with a valid UUID for testing
		CollectionID: testCollection.ID,     // Replace with a valid UUID for testing
		Embedding:    onesEmbedding,         // Replace with valid embedding data
		Limit:        10,                    // Replace with the desired limit for testing
		Threshold:    0.5,                   // Replace with the desired threshold for testing
	}

	// Insert test data into the database
	testData := []Document{
		{
			UserID:       testCollection.UserID,
			CollectionID: testCollection.ID,
			Filename:     "test1.pdf",
			Path:         "/path/to/test1.pdf",
			Pages: []*PageEmbedding{
				{
					Page:      0,
					Text:      "This is a test document",
					Embedding: onesEmbedding,
				},
			},
		},
		{
			UserID:       testCollection.UserID,
			CollectionID: testCollection.ID,
			Filename:     "test2.pdf",
			Path:         "/path/to/test2.pdf",
			Pages: []*PageEmbedding{
				{
					Page:      0,
					Text:      "This is a test2 document",
					Embedding: otherEmbedding,
				},
			},
		},
	}

	var ids []uuid.UUID

	for _, doc := range testData {
		docId, err := testClient.UpsertDocument(context.Background(), doc)
		if err != nil {
			t.Errorf("Failed to insert test data: %v", err)
			return
		}

		ids = append(ids, docId)
	}

	// Call the Search function
	results, err := testClient.Search(context.Background(), query)

	// Assert that the results are as expected
	expectedResults := []*SearchResult{
		// Define expected SearchResult objects based on your test data
		{
			DocumentID: ids[0],
			Filename:   testData[0].Filename,
			Pages: []*Page{
				{
					Page:  0,
					Text:  "This is a test document",
					Score: 1.0,
				},
			},
		},
	}

	// Check for errors
	if err != nil {
		t.Errorf("Search returned an error: %v", err)
		return
	}

	// Compare the actual results with expected results
	if len(results) != len(expectedResults) {
		t.Errorf("Expected %d results, but got %d", len(expectedResults), len(results))
		return
	}

	for i, actual := range results {
		expected := expectedResults[i]
		if actual.DocumentID != expected.DocumentID {
			t.Errorf("Expected DocumentID %s, but got %s", expected.DocumentID, actual.DocumentID)
		}

		// Perform similar comparisons for other fields in SearchResult

		// You can also check that the scores are within an acceptable range
		for j, page := range actual.Pages {
			if page.Score < query.Threshold {
				t.Errorf("Page %d has a score below the threshold", j)
			}
		}
	}
}
