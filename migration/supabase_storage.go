package migration

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	storage_go "github.com/supabase-community/storage-go"
	"log"
	"os"
)

type Supabase struct {
	Storage *storage_go.Client
	DB      *pgxpool.Pool
}

func InitSupabase(ctx context.Context) *Supabase {
	supaStorage := storage_go.NewClient(
		os.Getenv("SUPABASE_URL")+"/storage/v1",
		os.Getenv("SUPABASE_STORAGE_TOKEN"),
		nil)

	connection := os.Getenv("SUPABASE_DB")
	con, err := pgxpool.New(ctx, connection)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &Supabase{
		Storage: supaStorage,
		DB:      con,
	}
}

func (supabase Supabase) Close() {
	supabase.DB.Close()
}

func (supabase Supabase) StorageFiles(ctx context.Context) []string {

	// Read from storage schema table
	rows, err := supabase.DB.Query(ctx, "SELECT name FROM storage.objects")
	if err != nil {
		log.Fatalf("did not query: %v", err)
	}
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var path string
		err := rows.Scan(&path)
		if err != nil {
			log.Fatalf("did not scan: %v", err)
		}

		paths = append(paths, path)
	}

	return paths
}

func (supabase Supabase) GetFile(bucketId, path string) ([]byte, error) {
	return supabase.Storage.DownloadFile(bucketId, path)
}
