package migration

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"log"
	"time"
)

func (migrator *Migrator) MigrateUsages() {
	ctx := context.Background()

	log.Printf("Migrating usages...")

	rows, err := migrator.Legacy.Query(ctx, `SELECT 
		id, user_id, created_at, model, input_tokens, output_tokens 
		FROM model_usages`)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for rows.Next() {
		var (
			id      uuid.UUID
			userId  string
			date    time.Time
			modelId string
			input   uint32
			output  uint32
		)

		err = rows.Scan(&id, &userId, &date, &modelId, &input, &output)
		if err != nil {
			log.Fatal(err)
		}

		err = migrator.Next.InsertModelUsage(ctx, &datastore.ModelUsage{
			Id:           id,
			UserId:       userId,
			Timestamp:    date,
			ModelId:      modelId,
			InputTokens:  input,
			OutputTokens: output,
		})
		count++
	}

	log.Printf("Migrated %d usages", count)
}
