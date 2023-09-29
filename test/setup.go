package test

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/collections"
	"github.com/pzierahn/brainboost/documents"
	dbsetup "github.com/pzierahn/brainboost/setup"
	"github.com/sashabaranov/go-openai"
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
	documents   *documents.Service
	account     *account.Service
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

	acc := account.FromConfig(&account.Config{
		Auth: supabaseAuth,
		DB:   db,
	})

	gpt := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	return Setup{
		SupabaseUrl: supabaseUrl,
		Token:       token,
		db:          db,
		storage:     storage,
		collections: collections.NewServer(supabaseAuth, db, storage),
		account:     acc,
		documents: documents.FromConfig(&documents.Config{
			Auth:    supabaseAuth,
			Account: acc,
			DB:      db,
			GPT:     gpt,
			Storage: storage,
		}),
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
