package datastore

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	Content []*DocumentChunk `bson:"content,omitempty"`
}

type DocumentChunk struct {
	// ID of the document chunk
	Id uuid.UUID `bson:"id,omitempty"`

	// Content of the document chunk
	Text string `bson:"text,omitempty"`

	// Position of the document chunk
	Position uint32 `bson:"position,omitempty"`
}

// InsertDocument stores a document in the database.
func (service *Service) InsertDocument(ctx context.Context, document *Document) error {
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

// GetDocumentMeta retrieves the metadata of the documents from the database.
func (service *Service) GetDocumentMeta(ctx context.Context, userId string, ids ...uuid.UUID) ([]Document, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionDokuments)

	opts := &options.FindOptions{
		Projection: bson.M{
			"_id":    1,
			"name":   1,
			"type":   1,
			"source": 1,
		},
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
		"user_id": userId,
	}

	cursor, err := coll.Find(ctx, filter, opts)
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

type DocumentFilter struct {
	UserId       string
	CollectionId uuid.UUID
	Query        string
}

// ListDocuments retrieves all documents from the database without content.
func (service *Service) ListDocuments(ctx context.Context, query DocumentFilter) ([]Document, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionDokuments)

	filter := bson.M{
		"user_id":       query.UserId,
		"collection_id": query.CollectionId,
		"name": bson.M{
			"$regex": query.Query,
		},
	}

	opts := &options.FindOptions{
		Projection: bson.M{
			"_id":    1,
			"name":   1,
			"type":   1,
			"source": 1,
		},
	}

	cursor, err := coll.Find(ctx, filter, opts)
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

// RenameDocument renames a document in the database.
func (service *Service) RenameDocument(ctx context.Context, userId string, id uuid.UUID, name string) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionDokuments)

	_, err := coll.UpdateOne(ctx, bson.M{
		"_id":     id,
		"user_id": userId,
	}, bson.M{
		"$set": bson.M{
			"name": name,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteDocument deletes a document from the database.
func (service *Service) DeleteDocument(ctx context.Context, userId string, id uuid.UUID) error {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionDokuments)

	_, err := coll.DeleteOne(ctx, bson.M{
		"_id":     id,
		"user_id": userId,
	})
	if err != nil {
		return err
	}

	return nil
}

// FindRequest defines the input for finding a document ID
type FindRequest struct {
	UserId       string
	CollectionId uuid.UUID
	Name         string
}

// FindDocumentId finds a document ID for a given name
func (service *Service) FindDocumentId(ctx context.Context, req FindRequest) (uuid.UUID, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionDokuments)

	var document Document
	err := coll.FindOne(ctx, bson.M{
		"user_id":       req.UserId,
		"collection_id": req.CollectionId,
		"name":          req.Name,
	}).Decode(&document)
	if err != nil {
		return uuid.Nil, err
	}

	return document.Id, nil
}

func (service *Service) GetDocumentName(ctx context.Context, userId string, documentId uuid.UUID) (string, error) {
	coll := service.mongo.Database(DatabaseName).Collection(CollectionDokuments)

	opts := &options.FindOneOptions{
		Projection: bson.M{
			"name": 1,
		},
	}

	filter := bson.M{
		"user_id": userId,
		"_id":     documentId,
	}

	var document Document
	err := coll.FindOne(ctx, filter, opts).Decode(&document)
	if err != nil {
		return "", err
	}

	return document.Name, nil
}
