package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/chat"
	"github.com/pzierahn/brainboost/collections"
	"github.com/pzierahn/brainboost/documents"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/pzierahn/brainboost/setup"
	"github.com/sashabaranov/go-openai"
	storagego "github.com/supabase-community/storage-go"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
}

func main() {

	ctx := context.Background()

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

	jwtSec := os.Getenv("SUPABASE_JWT_SECRET")
	supabaseAuth := auth.WithSupabase(jwtSec)

	grpcServer := grpc.NewServer()

	collectionServer := collections.NewServer(supabaseAuth, db, storage)
	pb.RegisterCollectionServiceServer(grpcServer, collectionServer)

	accountService := account.FromConfig(&account.Config{
		Auth: supabaseAuth,
		DB:   db,
	})
	pb.RegisterAccountServiceServer(grpcServer, accountService)

	docsService := documents.FromConfig(&documents.Config{
		Auth:    supabaseAuth,
		Account: accountService,
		DB:      db,
		GPT:     gpt,
		Storage: storage,
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
