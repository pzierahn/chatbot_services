package main

import (
	"github.com/pzierahn/chatbot_services/chat"
	pb "github.com/pzierahn/chatbot_services/proto"
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

	grpcServer := grpc.NewServer()
	chatService, err := chat.New()
	if err != nil {
		log.Fatalf("failed to create chat service: %v", err)
	}
	pb.RegisterChatServiceServer(grpcServer, chatService)

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
