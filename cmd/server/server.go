package main

import (
	"flag"
	"github.com/pzierahn/braingain/braingain"
	"github.com/pzierahn/braingain/database"
	pb "github.com/pzierahn/braingain/proto"
	"github.com/pzierahn/braingain/server"
	"github.com/sashabaranov/go-openai"
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

	db, err := database.Connect("localhost:6334")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() { _ = db.Close() }()

	token := os.Getenv("OPENAI_API_KEY")
	gpt := openai.NewClient(token)

	chat := braingain.NewChat(db, gpt)

	doormanServer := server.NewServer(chat)
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
