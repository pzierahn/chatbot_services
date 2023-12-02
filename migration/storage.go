package migration

import (
	"context"
	firebase "firebase.google.com/go"
	"log"
	"path/filepath"
	"strings"
)

func Storage(supa *Supabase, app *firebase.App) {
	userMapping := GetUserIdMapping()

	ctx := context.Background()
	firebaseStorage, err := app.Storage(ctx)
	if err != nil {
		log.Fatalf("did not get storage: %v", err)
	}

	bucket, err := firebaseStorage.Bucket("brainboost-399710.appspot.com")
	if err != nil {
		log.Fatalf("did not get bucket: %v", err)
	}

	for _, path := range supa.StorageFiles(ctx) {
		byt, err := supa.GetFile("documents", path)
		if err != nil {
			log.Fatalf("did not download: %v", err)
		}

		pathParts := strings.Split(path, "/")
		googleId, ok := userMapping[pathParts[0]]
		if !ok {
			log.Fatalf("did not find user: %v", pathParts[0])
		}

		objPath := filepath.Join("documents", googleId, pathParts[1], pathParts[2])

		log.Printf("%s: %d", path, len(byt))

		// Write to firebase storage
		obj := bucket.Object(objPath)
		wrt := obj.NewWriter(ctx)
		_, err = wrt.Write(byt)
		if err != nil {
			log.Fatalf("did not write: %v", err)
		}

		err = wrt.Close()
		if err != nil {
			log.Fatalf("did not close: %v", err)
		}
	}
}
