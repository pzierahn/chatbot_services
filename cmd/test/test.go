package main

import (
	"github.com/pzierahn/chatbot_services/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

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

	tester := test.NewTester(conn)
	//tester.TestStartThread()
	//tester.TestThreadMessages()
	//tester.TestGetThread()
	//tester.TestListThreadIDs()
	//tester.TestDeleteThread()
	//tester.TestDeleteThreadMessage()
	tester.TestAccount()
}
