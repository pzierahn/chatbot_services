package test

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"time"
)

type Tester struct {
	db          *pgxpool.Pool
	chat        pb.ChatServiceClient
	collections pb.CollectionServiceClient
	account     pb.AccountServiceClient
	documents   pb.DocumentServiceClient
}

func (test Tester) createUser() string {
	uid := uuid.NewString()

	_, err := test.db.Exec(context.Background(),
		`INSERT INTO payments (user_id, amount)
			VALUES ($1, 1000)`, uid)

	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}

	return uid
}

func (test Tester) runTest(name string, testFunc func(ctx context.Context) error) {
	ctx := context.Background()

	uid := test.createUser()
	ctx = metadata.AppendToOutgoingContext(ctx, "User-Id", uid)

	start := time.Now()
	err := testFunc(ctx)
	execTime := time.Since(start)

	if err != nil {
		log.Printf("[%v] FAIL: %v", name, err)
	} else {
		log.Printf("[%v] PASS: %v", name, execTime)
	}
}

func (test Tester) expectError(name string, testFunc func(ctx context.Context) error) {
	ctx := context.Background()

	uid := test.createUser()
	ctx = metadata.AppendToOutgoingContext(ctx, "User-Id", uid)

	start := time.Now()
	err := testFunc(ctx)
	execTime := time.Since(start)

	if err == nil {
		log.Printf("[%v] FAIL: expected error", name)
	} else {
		log.Printf("[%v] PASS: %v", name, execTime)
	}
}

func NewTester(conn *grpc.ClientConn) *Tester {

	ctx := context.Background()
	addr := os.Getenv("CHATBOT_DB_TEST")
	db, err := pgxpool.New(ctx, addr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &Tester{
		db:          db,
		chat:        pb.NewChatServiceClient(conn),
		collections: pb.NewCollectionServiceClient(conn),
		account:     pb.NewAccountServiceClient(conn),
		documents:   pb.NewDocumentServiceClient(conn),
	}
}
