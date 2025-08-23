package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const credentialsFile = "service_account.json"
const bucketName = "brainboost-399710.appspot.com"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func initFirebase(ctx context.Context) *firebase.App {
	var opts []option.ClientOption
	if _, err := os.Stat(credentialsFile); err == nil {
		opts = append(opts, option.WithCredentialsFile(credentialsFile))
	}
	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		log.Fatalf("failed to create firebase app: %v", err)
	}
	return app
}

func initBucket(ctx context.Context, app *firebase.App) *storage.BucketHandle {
	fs, err := app.Storage(ctx)
	if err != nil {
		log.Fatalf("failed to create firebase storage client: %v", err)
	}
	bucket, err := fs.Bucket(bucketName)
	if err != nil {
		log.Fatalf("failed to get bucket %q: %v", bucketName, err)
	}
	return bucket
}

func ensureDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0o755)
}

func downloadObject(ctx context.Context, bucket *storage.BucketHandle, name, destRoot string) error {
	// Skip empty or directory-like placeholders
	if name == "" || strings.HasSuffix(name, "/") {
		return nil
	}

	obj := bucket.Object(name)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("open reader for %s: %w", name, err)
	}
	defer r.Close()

	localPath := filepath.Join(destRoot, filepath.FromSlash(name))
	if err := ensureDir(localPath); err != nil {
		return fmt.Errorf("create dir for %s: %w", localPath, err)
	}

	tmpPath := localPath + ".part"
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create file %s: %w", tmpPath, err)
	}

	if _, err = io.Copy(f, r); err != nil {
		_ = f.Close()
		return fmt.Errorf("write file %s: %w", tmpPath, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close file %s: %w", tmpPath, err)
	}

	if err := os.Rename(tmpPath, localPath); err != nil {
		return fmt.Errorf("finalize file %s: %w", localPath, err)
	}

	return nil
}

func main() {
	// Flags
	destDir := flag.String("out", "data", "Destination directory for downloads")
	prefix := flag.String("prefix", "", "Only download objects with this prefix")
	workers := flag.Int("workers", 4, "Number of concurrent downloads")
	flag.Parse()

	ctx := context.Background()
	app := initFirebase(ctx)
	bucket := initBucket(ctx, app)

	if err := os.MkdirAll(*destDir, 0o755); err != nil {
		log.Fatalf("failed to create dest dir: %v", err)
	}

	// List objects
	query := &storage.Query{}
	if *prefix != "" {
		query.Prefix = *prefix
	}
	it := bucket.Objects(ctx, query)

	// Worker pool
	type job struct{ name string }
	jobs := make(chan job, 64)
	errs := make(chan error, *workers)

	for i := 0; i < *workers; i++ {
		go func(id int) {
			for j := range jobs {
				start := time.Now()
				if err := downloadObject(ctx, bucket, j.name, *destDir); err != nil {
					errs <- fmt.Errorf("%s: %w", j.name, err)
					continue
				}
				log.Printf("[%d] downloaded %s in %s", id, j.name, time.Since(start).Truncate(time.Millisecond))
			}
			errs <- nil
		}(i + 1)
	}

	count := 0
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("listing objects failed: %v", err)
		}
		jobs <- job{name: attrs.Name}
		count++
	}
	close(jobs)

	// Wait for workers
	var hadErr bool
	for i := 0; i < *workers; i++ {
		if err := <-errs; err != nil {
			log.Printf("error: %v", err)
			hadErr = true
		}
	}

	if hadErr {
		log.Fatalf("completed with errors")
	}
	log.Printf("done. downloaded %d object(s) to %s", count, *destDir)
}
