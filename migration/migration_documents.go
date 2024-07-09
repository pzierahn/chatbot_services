package migration

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/search"
	"log"
	"sync"
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

func (migrator *Migrator) MigrateDocumentToSearch(index search.Index) {
	ctx := context.Background()

	log.Printf("Migrating documents")

	//
	// Get all documents
	//

	// Query all document and user ids
	rows, err := migrator.Legacy.Query(ctx, "SELECT id, user_id, collection_id FROM documents")
	if err != nil {
		log.Fatalf("Query documents: %v", err)
	}

	upserts := make(chan int, 3)
	defer close(upserts)

	var wg sync.WaitGroup

	go func() {
		processedFragments := 0
		processedDocuments := 0

		for upsert := range upserts {
			processedFragments += upsert
			processedDocuments++
			log.Printf("Processed: docs=%d fragments=%d", processedDocuments, processedFragments)
		}
	}()

	// Iterate over all documents
	for rows.Next() {
		var (
			docId        uuid.UUID
			userId       string
			collectionId uuid.UUID
		)

		err = rows.Scan(&docId, &userId, &collectionId)
		if err != nil {
			log.Fatalf("Scan document: %v", err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			doc, err := migrator.Next.GetDocument(ctx, userId, docId)
			if err != nil {
				log.Fatalf("Get document: %v", err)
			}

			fragments := make([]*search.Fragment, len(doc.Content))
			for idx, chunk := range doc.Content {
				fragments[idx] = &search.Fragment{
					Id:           chunk.Id.String(),
					Text:         chunk.Text,
					UserId:       userId,
					DocumentId:   doc.Id.String(),
					CollectionId: collectionId.String(),
					Position:     chunk.Position,
				}
			}

			_, err = index.Upsert(ctx, fragments)
			if err != nil {
				log.Fatalf("Upsert fragments: %v", err)
			}

			upserts <- len(fragments)
		}()
	}

	wg.Wait()
}
