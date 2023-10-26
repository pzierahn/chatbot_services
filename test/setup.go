package test

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/chat"
	"github.com/pzierahn/brainboost/collections"
	"github.com/pzierahn/brainboost/documents"
	pb "github.com/pzierahn/brainboost/proto"
	dbsetup "github.com/pzierahn/brainboost/setup"
	"github.com/sashabaranov/go-openai"
	storagego "github.com/supabase-community/storage-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
)

type Setup struct {
	SupabaseUrl string
	Token       string
	lis         net.Listener
	db          *pgxpool.Pool
	storage     *storagego.Client
	collections pb.CollectionServiceClient
	documents   pb.DocumentServiceClient
	chat        pb.ChatServiceClient
	account     pb.AccountServiceClient
	Report      *Report
}

const bucket = "documents"

func NewTestSetup() Setup {
	supabaseUrl := os.Getenv("API_EXTERNAL_URL")
	token := os.Getenv("SERVICE_ROLE_KEY")
	postgresDB := "postgres://postgres:your-super-secret-and-long-postgres-password@localhost:5432/postgres"

	storage := storagego.NewClient(supabaseUrl+"/storage/v1", token, nil)
	_, _ = storage.CreateBucket(bucket, storagego.BucketOptions{
		Public:        false,
		FileSizeLimit: "50mb",
	})

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

	docServer := documents.FromConfig(&documents.Config{
		Auth:    supabaseAuth,
		Account: acc,
		DB:      db,
		GPT:     gpt,
		Storage: storage,
	})

	collectionsServer :=
		collections.NewServer(supabaseAuth, db, storage)

	grpcServer := grpc.NewServer()
	pb.RegisterAccountServiceServer(grpcServer, acc)
	pb.RegisterDocumentServiceServer(grpcServer, docServer)
	pb.RegisterCollectionServiceServer(grpcServer, collectionsServer)
	pb.RegisterChatServiceServer(grpcServer, chat.FromConfig(&chat.Config{
		DB:              db,
		GPT:             gpt,
		AuthService:     supabaseAuth,
		DocumentService: docServer,
		AccountService:  acc,
	}))

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		_ = grpcServer.Serve(lis)
	}()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return Setup{
		SupabaseUrl: supabaseUrl,
		Token:       token,
		db:          db,
		storage:     storage,
		collections: pb.NewCollectionServiceClient(conn),
		account:     pb.NewAccountServiceClient(conn),
		documents:   pb.NewDocumentServiceClient(conn),
		chat:        pb.NewChatServiceClient(conn),
		lis:         lis,
		Report:      &Report{},
	}
}

func (setup *Setup) Close() {

	err := setup.lis.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = dbsetup.DropTables(context.Background(), setup.db)
	if err != nil {
		log.Fatal(err)
	}

	setup.db.Close()

	_, emErr := setup.storage.EmptyBucket(bucket)
	if emErr != nil {
		log.Fatal(emErr)
	}

	_, delError := setup.storage.DeleteBucket(bucket)
	if delError != nil {
		log.Fatal(delError)
	}

	log.Printf("Deleted bucket: %v", bucket)
}
