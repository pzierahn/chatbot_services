package collections

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/auth"
	pb "github.com/pzierahn/brainboost/proto"
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

type mockCollections struct {
	userId      uuid.UUID
	collections []*pb.Collection
}

func setupCollections(t *testing.T) []mockCollections {
	ctx := context.Background()

	users := []struct {
		id          uuid.UUID
		collections []*pb.Collection
	}{
		{
			id: uuid.New(),
			collections: []*pb.Collection{
				{
					Name: "Collection 1",
				},
				{
					Name: "Collection 2",
				},
			},
		},
		{
			id: uuid.New(),
			collections: []*pb.Collection{
				{
					Name: "Data 1",
				},
				{
					Name: "Data 2",
				},
			},
		},
	}

	var mocks []mockCollections

	// Create users and collections
	for _, user := range users {
		service := &Service{
			auth: auth.WithUser(user.id),
			db:   testConn,
		}

		mock := mockCollections{
			userId: user.id,
		}

		for _, req := range user.collections {
			collection, err := service.Create(ctx, req)
			if err != nil {
				t.Fatal(err)
			}

			mock.collections = append(mock.collections, collection)
		}

		mocks = append(mocks, mock)
	}

	return mocks
}
