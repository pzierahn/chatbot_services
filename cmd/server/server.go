package main

import (
	"context"
	"crypto/tls"
	firebase "firebase.google.com/go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/chat"
	"github.com/pzierahn/brainboost/collections"
	"github.com/pzierahn/brainboost/documents"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/pzierahn/brainboost/setup"
	"github.com/sashabaranov/go-openai"
	storagego "github.com/supabase-community/storage-go"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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

	addr := os.Getenv("SUPABASE_DB")
	db, err := pgxpool.New(ctx, addr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer db.Close()

	err = setup.SetupTables(ctx, db)
	if err != nil {
		log.Fatalf("failed to setup tables: %v", err)
	}

	token := os.Getenv("OPENAI_API_KEY")
	gpt := openai.NewClient(token)

	storage := storagego.NewClient(
		os.Getenv("SUPABASE_URL")+"/storage/v1",
		os.Getenv("SUPABASE_STORAGE_TOKEN"),
		nil)

	supabaseAuth, err := auth.WithFirebase(ctx, app)
	if err != nil {
		log.Fatalf("failed to create auth service: %v", err)
	}

	grpcServer := grpc.NewServer()

	collectionServer := collections.NewServer(supabaseAuth, db, storage)
	pb.RegisterCollectionServiceServer(grpcServer, collectionServer)

	accountService := account.FromConfig(&account.Config{
		Auth: supabaseAuth,
		DB:   db,
	})
	pb.RegisterAccountServiceServer(grpcServer, accountService)

	config := &tls.Config{}

	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", os.Getenv("PINECONE_KEY"))
	target := os.Getenv("PINECONE_URL")

	log.Printf("connecting to %v", target)

	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithAuthority(target),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	pineconeClient := pinecone_grpc.NewVectorServiceClient(conn)

	docsService := documents.FromConfig(&documents.Config{
		Auth:     supabaseAuth,
		Account:  accountService,
		DB:       db,
		GPT:      gpt,
		Storage:  storage,
		Pinecone: pineconeClient,
	})
	pb.RegisterDocumentServiceServer(grpcServer, docsService)

	chatServer := chat.FromConfig(&chat.Config{
		DB:              db,
		GPT:             gpt,
		DocumentService: docsService,
		AccountService:  accountService,
		AuthService:     supabaseAuth,
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
