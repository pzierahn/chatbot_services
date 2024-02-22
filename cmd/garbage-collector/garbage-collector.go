package main

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/chatbot_services/utils"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"
)

const credentialsFile = "service_account.json"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	var opts []option.ClientOption
	if _, err := os.Stat(credentialsFile); err == nil {
		serviceAccount := option.WithCredentialsFile(credentialsFile)
		opts = append(opts, serviceAccount)
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		log.Fatalf("failed to create firebase app: %v", err)
	}

	firebaseStorage, err := app.Storage(ctx)
	if err != nil {
		log.Fatalf("failed to create firebase storage client: %v", err)
	}

	bucket, err := firebaseStorage.Bucket("brainboost-399710.appspot.com")
	if err != nil {
		log.Fatalf("did not get bucket: %v", err)
	}

	// List all objects in the bucket
	iter := bucket.Objects(ctx, nil)
	if err != nil {
		log.Fatalf("failed to list objects: %v", err)
	}

	addr := os.Getenv("CHATBOT_DB")
	db, err := pgxpool.New(ctx, addr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer db.Close()

	userIds := make(map[string][]string)
	var freeedSpace int64

	for {
		attrs, err := iter.Next()
		// Check if err is EOF
		if err != nil {
			break
		}

		parts := strings.Split(attrs.Name, "/")
		file := parts[len(parts)-1]

		documentId := file[:len(file)-4]
		userId := parts[1]

		// Check if the document exists in the database
		var exists bool
		err = db.QueryRow(
			ctx,
			`SELECT EXISTS(SELECT 1 FROM documents WHERE id = $1)`,
			documentId).Scan(&exists)
		if err != nil {
			log.Fatalf("failed to query: %v", err)
		}

		if !exists {
			userIds[userId] = append(userIds[userId], documentId)

			log.Printf("--> delete %v %v", attrs.Name, attrs.Size)

			// Delete the file
			err = bucket.Object(attrs.Name).Delete(ctx)
			if err != nil {
				log.Fatalf("failed to delete object: %v", err)
			}

			freeedSpace += attrs.Size
		}
	}

	log.Printf("userIds: %s", utils.Prettify(userIds))
	log.Printf("freed space: %d", freeedSpace)
}
