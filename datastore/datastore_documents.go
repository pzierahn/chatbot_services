package datastore

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	DocumentTypePDF = "pdf"
	DocumentTypeWeb = "web"
)

type Document struct {
	// ID of the document
	Id uuid.UUID `bson:"_id,omitempty"`

	// User ID
	UserId string `bson:"user_id,omitempty"`

	// Collection ID
	CollectionId uuid.UUID `bson:"collection_id,omitempty"`

	// Name of the document
	Name string `bson:"name,omitempty"`

	// Type of the document
	Type string `bson:"type,omitempty"`

	// Source can be a URL or a file path
	Source string `bson:"source,omitempty"`

	// Data chunks
	Content []DocumentChunk `bson:"content,omitempty"`
}

type DocumentChunk struct {
	// ID of the document chunk
	Id uuid.UUID `bson:"id,omitempty"`

	// Content of the document chunk
	Text string `bson:"text,omitempty"`

	// Position of the document chunk
	Position int `bson:"position,omitempty"`
}

// StoreDocument stores a document in the database.
func (service *Service) StoreDocument(ctx context.Context, document *Document) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionDokuments)

	_, err := coll.InsertOne(ctx, document)
	if err != nil {
		return err
	}

	return nil
}

// GetDocument retrieves a document from the database.
func (service *Service) GetDocument(ctx context.Context, userId string, id uuid.UUID) (*Document, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionDokuments)

	var document Document
	err := coll.FindOne(ctx, bson.M{
		"_id":     id,
		"user_id": userId,
	}).Decode(&document)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

// GetDocuments retrieves all documents from the database.
func (service *Service) GetDocuments(ctx context.Context, userId string, ids ...uuid.UUID) ([]Document, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionDokuments)

	cursor, err := coll.Find(ctx, bson.M{
		"_id": bson.M{
			"$in": ids,
		},
		"user_id": userId,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(ctx) }()

	var documents []Document
	err = cursor.All(ctx, &documents)
	if err != nil {
		return nil, err
	}

	return documents, nil
}
