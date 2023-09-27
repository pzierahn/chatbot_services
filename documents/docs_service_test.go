package documents

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/setup"
	"os"
	"testing"
)

var testService *Service

func setupTestService(t *testing.T) {
	testConnection := os.Getenv("TEST_DATABASE")

	ctx := context.Background()

	conn, err := pgxpool.New(ctx, testConnection)
	if err != nil {
		t.Fatal(err)
	}

	userAuth := auth.WithUser(uuid.New())
	user := account.NewServer(userAuth, conn)

	testService = FromConfig(&Config{
		Auth:    userAuth,
		Account: user,
		Db:      conn,
		Gpt:     nil,
		Storage: nil,
	})
	err = setup.SetupTables(ctx, testService.db)
	if err != nil {
		t.Fatal(err)
	}
}

func teardownTestService(t *testing.T) {
	ctx := context.Background()
	err := setup.DropTables(ctx, testService.db)
	if err != nil {
		t.Fatal(err)
	}
}
