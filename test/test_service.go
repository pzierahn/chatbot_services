package test

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/collections"
	dbsetup "github.com/pzierahn/brainboost/setup"
	storagego "github.com/supabase-community/storage-go"
	"log"
	"os"
)

type Setup struct {
	SupabaseUrl string
	Token       string
	db          *pgxpool.Pool
	storage     *storagego.Client
	collections *collections.Service
}

func NewTestSetup() Setup {
	supabaseUrl := os.Getenv("API_EXTERNAL_URL")
	token := os.Getenv("SERVICE_ROLE_KEY")
	postgresDB := "postgres://postgres:your-super-secret-and-long-postgres-password@localhost:5432/postgres"

	storage := storagego.NewClient(supabaseUrl+"/storage/v1", token, nil)

	ctx := context.Background()

	db, err := pgxpool.New(ctx, postgresDB)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	err = dbsetup.SetupTables(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	supabaseAuth := auth.WithSupabase()

	return Setup{
		SupabaseUrl: supabaseUrl,
		Token:       token,
		db:          db,
		storage:     storage,
		collections: collections.NewServer(supabaseAuth, db, storage),
	}
}

func (setup *Setup) Close() {

	err := dbsetup.DropTables(context.Background(), setup.db)
	if err != nil {
		log.Fatal(err)
	}

	setup.db.Close()

	_, errr := setup.storage.DeleteBucket("documents")
	if errr.Error != "" {
		log.Fatal(errr)
	}
}
