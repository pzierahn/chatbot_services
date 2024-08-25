package migration

import (
	"context"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/search"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func (migrator *Migrator) MigrateVectorDB() {
	ctx := context.Background()

	log.Printf("Migrating documents...")

	collection := migrator.Database.Database(datastore.DatabaseName).Collection(datastore.CollectionDokuments)
	cur, err := collection.Find(ctx, &bson.M{})
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer func() { _ = cur.Close(ctx) }()

	var totalUsage uint32
	var idx uint32

	for cur.Next(ctx) {
		var doc datastore.Document
		err := cur.Decode(&doc)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}

		log.Printf("[%3d] Migrating document %s (%d)", idx, doc.Id, len(doc.Content))
		var fragments []*search.Fragment

		for _, chunk := range doc.Content {
			if chunk.Text == "" {
				continue
			}

			fragments = append(fragments, &search.Fragment{
				Id:           chunk.Id.String(),
				Text:         chunk.Text,
				UserId:       doc.UserId,
				DocumentId:   doc.Id.String(),
				CollectionId: doc.CollectionId.String(),
				Position:     chunk.Position,
			})
		}

		// Upsert fragments in chunks of 100 to avoid too large requests.
		for start := 0; start < len(fragments); start += 100 {
			end := min(start+100, len(fragments))
			usage, err := migrator.Search.Upsert(ctx, fragments[start:end])
			if err != nil {
				log.Fatalf("Error: %s", err)
			}

			totalUsage += usage.Tokens
		}

		idx++
	}

	log.Printf("Total tokens: %d", totalUsage)
	log.Printf("Migration done")
}
