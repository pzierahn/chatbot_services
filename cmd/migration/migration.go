package main

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/jackc/pgx/v5/pgxpool"
	storage_go "github.com/supabase-community/storage-go"
	"google.golang.org/api/option"
	"log"
	"os"
)

var app *firebase.App

func init() {
	opt := option.WithCredentialsFile("service_account.json")

	var err error
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	connection := os.Getenv("SUPABASE_DB")
	db, err := pgxpool.New(ctx, connection)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer db.Close()

	err = db.Ping(ctx)
	if err != nil {
		log.Fatalf("did not ping: %v", err)
	}

	log.Printf("connected")

	// Read from storage schema table
	rows, err := db.Query(ctx, "SELECT name FROM storage.objects")
	if err != nil {
		log.Fatalf("did not query: %v", err)
	}
	defer rows.Close()

	supaStorage := storage_go.NewClient(
		os.Getenv("SUPABASE_URL")+"/storage/v1",
		os.Getenv("SUPABASE_STORAGE_TOKEN"),
		nil)

	var paths []string
	for rows.Next() {
		var path string
		err := rows.Scan(&path)
		if err != nil {
			log.Fatalf("did not scan: %v", err)
		}

		// log.Printf("%s", path)
		paths = append(paths, path)
	}

	firebaseStorage, err := app.Storage(ctx)
	if err != nil {
		log.Fatalf("did not get storage: %v", err)
	}

	bucket, err := firebaseStorage.Bucket("brainboost-399710.appspot.com")
	if err != nil {
		log.Fatalf("did not get bucket: %v", err)
	}

	for _, path := range paths {
		byt, err := supaStorage.DownloadFile("documents", path)
		if err != nil {
			log.Fatalf("did not download: %v", err)
		}

		log.Printf("Download size: %d", len(byt))

		// Write to firebase storage
		obj := bucket.Object(path)
		wrt := obj.NewWriter(ctx)
		n, err := wrt.Write(byt)
		if err != nil {
			log.Fatalf("did not write: %v", err)
		}
		err = wrt.Close()
		if err != nil {
			log.Fatalf("did not close: %v", err)
		}

		log.Printf("Write size: %d", n)

		break
	}
}
