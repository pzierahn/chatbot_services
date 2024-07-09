package migration

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"log"
)

type Webpage struct {
	Url   string `json:"url,omitempty"`
	Title string `json:"title,omitempty"`
}

type File struct {
	Path     string `json:"path,omitempty"`
	Filename string `json:"filename,omitempty"`
}

type DocumentMeta struct {
	Webpage *Webpage `json:"webpage,omitempty"`
	File    *File    `json:"file,omitempty"`
}

type document struct {
	ID           uuid.UUID    `json:"id,omitempty"`
	UserID       string       `json:"user_id,omitempty"`
	CollectionID uuid.UUID    `json:"collection_id,omitempty"`
	Metadata     DocumentMeta `json:"metadata,omitempty"`
}

func (migrator *Migrator) MigrateDocuments() {
	ctx := context.Background()

	log.Printf("Migrating documents")

	//
	// Get all documents
	//

	documents := make(map[uuid.UUID]document)

	rows, err := migrator.Legacy.Query(ctx, "SELECT id, user_id, collection_id, metadata FROM documents")
	if err != nil {
		log.Fatalf("Query collections: %v", err)
	}
	for rows.Next() {
		var doc document

		err = rows.Scan(&doc.ID, &doc.UserID, &doc.CollectionID, &doc.Metadata)
		if err != nil {
			log.Fatalf("Scan collection: %v", err)
		}

		documents[doc.ID] = doc
	}

	//
	// Get all document chunks
	//

	docChunks := make(map[uuid.UUID][]*datastore.DocumentChunk)

	rows, err = migrator.Legacy.Query(ctx, "SELECT document_id, id, index, text FROM document_chunks")
	if err != nil {
		log.Fatalf("Query document_chunks: %v", err)
	}
	for rows.Next() {
		var docId uuid.UUID
		var chunk datastore.DocumentChunk

		err = rows.Scan(&docId, &chunk.Id, &chunk.Position, &chunk.Text)
		if err != nil {
			log.Fatalf("Scan document_chunk: %v", err)
		}

		docChunks[docId] = append(docChunks[docId], &chunk)
	}

	//
	// Match chunks to documents and insert them into the new datastore
	//

	for docId, doc := range documents {
		chunks, ok := docChunks[docId]
		if !ok {
			log.Fatalf("Document chunks not found")
		}

		var name, docType, source string
		if doc.Metadata.File != nil {
			name = doc.Metadata.File.Filename
			docType = datastore.DocumentTypePDF
			source = doc.Metadata.File.Path
		}
		if doc.Metadata.Webpage != nil {
			name = doc.Metadata.Webpage.Title
			docType = datastore.DocumentTypeWeb
			source = doc.Metadata.Webpage.Url
		}

		nextDocument := &datastore.Document{
			Id:           doc.ID,
			UserId:       doc.UserID,
			CollectionId: doc.CollectionID,
			Name:         name,
			Type:         docType,
			Source:       source,
			Content:      chunks,
		}

		err = migrator.Next.InsertDocument(ctx, nextDocument)
		if err != nil {
			log.Fatalf("Insert document: %v", err)
		}
	}
}
