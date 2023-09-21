package main

import (
	"context"
	"flag"
	"github.com/pzierahn/brainboost/database"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/pzierahn/brainboost/server"
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
	db, err := database.Connect(ctx, os.Getenv("SUPABASE_DB"))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer db.Close()

	err = db.SetupTables(ctx)
	if err != nil {
		log.Fatalf("failed to setup tables: %v", err)
	}

	token := os.Getenv("OPENAI_API_KEY")
	gpt := openai.NewClient(token)

	storage := storagego.NewClient(
		os.Getenv("SUPABASE_URL")+"/storage/v1",
		os.Getenv("SUPABASE_STORAGE_TOKEN"),
		nil)

	doormanServer := server.NewServer(db, gpt, storage)
	grpcServer := grpc.NewServer()
	pb.RegisterBraingainServer(grpcServer, doormanServer)

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
