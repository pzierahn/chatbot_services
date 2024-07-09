package migration

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/llm"
	"log"
	"sort"
	"time"
)

type threadMessage struct {
	id           uuid.UUID
	userId       string
	threadId     uuid.UUID
	createdAt    time.Time
	prompt       string
	completion   string
	collectionId uuid.UUID
}

func (migrator *Migrator) MigrateThreads() {
	ctx := context.Background()

	log.Printf("Migrating threads...")

	rows, err := migrator.Legacy.Query(ctx, `SELECT
		tm.id AS thread_message_id,
		tm.user_id AS thread_message_user_id,
		tm.thread_id,
		tm.created_at AS thread_message_created_at,
		tm.prompt,
		tm.completion,
		t.collection_id
	FROM
		thread_messages tm
	JOIN
		threads t ON tm.thread_id = t.id;`)
	if err != nil {
		log.Fatalf("Query collections: %v", err)
	}
	defer rows.Close()

	threadMessages := make(map[uuid.UUID][]*threadMessage)
	for rows.Next() {
		var message threadMessage

		err := rows.Scan(
			&message.id,
			&message.userId,
			&message.threadId,
			&message.createdAt,
			&message.prompt,
			&message.completion,
			&message.collectionId,
		)
		if err != nil {
			log.Fatalf("Scan thread message: %v", err)
		}

		threadMessages[message.threadId] = append(threadMessages[message.threadId], &message)
	}

	for threadId, messages := range threadMessages {
		// Sort from oldest to newest
		sort.Slice(messages, func(i, j int) bool {
			return messages[i].createdAt.Before(messages[j].createdAt)
		})

		thread := &datastore.Thread{
			Id:           threadId,
			UserId:       messages[0].userId,
			CollectionId: messages[0].collectionId,
			Timestamp:    messages[0].createdAt,
		}

		for _, message := range messages {
			thread.Messages = append(thread.Messages, []*llm.Message{
				{
					Role:    llm.RoleUser,
					Content: message.prompt,
				},
				{
					Role:    llm.RoleAssistant,
					Content: message.completion,
				},
			}...)
		}

		err = migrator.Next.StoreThread(ctx, thread)
		if err != nil {
			log.Fatalf("Store thread: %v", err)
		}
	}

	log.Printf("Migrated %d threads.", len(threadMessages))
}
