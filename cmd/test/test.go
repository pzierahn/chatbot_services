package main

import (
	"encoding/json"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func prettify(obj interface{}) string {
	byt, _ := json.MarshalIndent(obj, "", "  ")
	return string(byt)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conn, err := grpc.Dial(
		"localhost:8869",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() { _ = conn.Close() }()

	service := pb.NewChatServiceClient(conn)

	tester := test.NewTester(service)
	tester.TestThreads()
}
