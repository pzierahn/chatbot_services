package main

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/chat"
	"github.com/pzierahn/brainboost/collections"
	"github.com/pzierahn/brainboost/documents"
	"github.com/pzierahn/brainboost/llm/openai"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/pzierahn/brainboost/setup"
	"github.com/pzierahn/brainboost/vectordb"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const credentialsFile = "service_account.json"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	ctx := context.Background()

	var opts []option.ClientOption

	if _, err := os.Stat(credentialsFile); err == nil {
		serviceAccount := option.WithCredentialsFile(credentialsFile)
		opts = append(opts, serviceAccount)
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		log.Fatalf("failed to create firebase app: %v", err)
	}

	firebaseStorage, err := app.Storage(ctx)
	if err != nil {
		log.Fatalf("failed to create firebase storage client: %v", err)
	}

	bucket, err := firebaseStorage.Bucket("brainboost-399710.appspot.com")
	if err != nil {
		log.Fatalf("did not get bucket: %v", err)
	}

	addr := os.Getenv("BRAINBOOST_COCKROACH_DB")
	db, err := pgxpool.New(ctx, addr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer db.Close()

	err = setup.SetupTables(ctx, db)
	if err != nil {
		log.Fatalf("failed to setup tables: %v", err)
	}

	embeddingService := openai.New()
	completionService := embeddingService

	authService, err := auth.WithFirebase(ctx, app)
	if err != nil {
		log.Fatalf("failed to create auth service: %v", err)
	}

	vecDB, err := vectordb.New()
	if err != nil {
		log.Fatalf("failed to create vector db: %v", err)
	}
	defer func() { _ = vecDB.Close() }()

	grpcServer := grpc.NewServer()
	collectionServer := collections.NewServer(authService, db, bucket, vecDB)
	pb.RegisterCollectionServiceServer(grpcServer, collectionServer)

	accountService := account.FromConfig(&account.Config{
		Auth: authService,
		DB:   db,
	})
	pb.RegisterAccountServiceServer(grpcServer, accountService)

	docsService := documents.FromConfig(&documents.Config{
		Auth:       authService,
		Account:    accountService,
		DB:         db,
		Embeddings: embeddingService,
		Storage:    bucket,
		VectorDB:   vecDB,
	})
	pb.RegisterDocumentServiceServer(grpcServer, docsService)

	chatServer := chat.FromConfig(&chat.Config{
		DB:              db,
		Completion:      completionService,
		DocumentService: docsService,
		AccountService:  accountService,
		AuthService:     authService,
	})
	pb.RegisterChatServiceServer(grpcServer, chatServer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9055"
		log.Printf("defaulting to port %s", port)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("starting server on %v", listener.Addr().String())
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
