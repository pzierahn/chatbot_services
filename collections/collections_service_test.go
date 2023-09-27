package collections

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/setup"
	"os"
	"testing"
)

var testService *Service
var testConn *pgxpool.Pool

func setupTestService(t *testing.T) {
	testConnection := os.Getenv("TEST_DATABASE")

	ctx := context.Background()

	conn, err := pgxpool.New(ctx, testConnection)
	if err != nil {
		t.Fatal(err)
	}

	userAuth := auth.WithUser(uuid.New())

	testService = NewServer(userAuth, conn, nil)
	err = setup.SetupTables(ctx, testService.db)
	if err != nil {
		t.Fatal(err)
	}

	testConn = conn
}

func teardownTestService(t *testing.T) {
	ctx := context.Background()
	err := setup.DropTables(ctx, testService.db)
	if err != nil {
		t.Fatal(err)
	}
}
