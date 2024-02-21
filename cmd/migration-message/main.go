package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

type ChatMessage struct {
	ID           string
	UserID       string
	CreatedAt    time.Time
	CollectionID string
	Prompt       string
	Completion   string
}

type migrationService struct {
	db *pgxpool.Pool
}

func (mig migrationService) getAllChatMessages() (messages []ChatMessage) {
	ctx := context.Background()
	rows, err := mig.db.Query(ctx, "SELECT * FROM chat_messages")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var chatMessage ChatMessage
		err := rows.Scan(
			&chatMessage.ID,
			&chatMessage.UserID,
			&chatMessage.CreatedAt,
			&chatMessage.CollectionID,
			&chatMessage.Prompt,
			&chatMessage.Completion,
		)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
		}

		messages = append(messages, chatMessage)
	}

	return messages
}

func (mig migrationService) migrateMessage(message ChatMessage) {
	ctx := context.Background()

	tx, err := mig.db.Begin(ctx)
	if err != nil {
		log.Fatalf("Begin failed: %v", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO threads (id, user_id, created_at, collection_id)
				VALUES ($1, $2, $3, $4)`,
		message.ID,
		message.UserID,
		message.CreatedAt,
		message.CollectionID)
	if err != nil {
		log.Fatalf("Insert failed: %v", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO thread_messages (user_id, thread_id, created_at, prompt, completion)
				VALUES ($1, $2, $3, $4, $5)`,
		message.UserID,
		message.ID,
		message.CreatedAt,
		message.Prompt,
		message.Completion)
	if err != nil {
		log.Fatalf("Insert failed: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatalf("Commit failed: %v", err)
	}
}

func (mig migrationService) migrateRefs(message ChatMessage) {
	ctx := context.Background()

	rows, err := mig.db.Query(ctx,
		`SELECT document_chunk_id
			FROM chat_message_references
			WHERE chat_message_id = $1`,
		message.ID)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var chunkId string
		err = rows.Scan(&chunkId)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
		}

		_, err = mig.db.Exec(ctx,
			`INSERT INTO thread_references (user_id, thread_id, document_chunk_id)
				VALUES ($1, $2, $3)`,
			message.UserID,
			message.ID,
			chunkId)
		if err != nil {
			log.Fatalf("Insert failed: %v", err)
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	addr := os.Getenv("CHATBOT_DB")
	db, err := pgxpool.New(ctx, addr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer db.Close()

	mig := migrationService{db: db}
	messages := mig.getAllChatMessages()
	log.Printf("Messages: %v", len(messages))

	for inx, message := range messages[1:] {
		log.Printf("Message[%v]: %v", inx, message.ID)
		mig.migrateMessage(message)
		mig.migrateRefs(message)
	}
}
