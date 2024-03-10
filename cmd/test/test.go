package main

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	meta := metadata.New(map[string]string{
		"User-Id": "j7jjxLD9rla2DrZoeUu3Tnft4812",
	})

	conn, err := grpc.Dial(
		"localhost:8869",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// Add header to all requests
		grpc.WithDefaultCallOptions(grpc.Header(&meta)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() { _ = conn.Close() }()

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, meta)

	table := &pb.NewTable{
		Name: "Test Table",
		//DocumentIds: []string{
		//	// Efficient_Byzantine_Fault-Tolerance.pdf
		//	"6436db5c-44e9-4468-b716-bbb701739d18",
		//	// Lamport - Time, Clocks, and the Ordering of Events in a Distributed System.pdf
		//	"828cbaa2-22d7-4659-b2b5-75ace9a43fbf",
		//	// Kademlia.pdf
		//	"bcd64239-16cf-4f38-80d2-1ff6ea3ccda5",
		//},
		//Columns: []string{
		//	"Extract the title and return only the title text, without any additional words or explanation",
		//	"Extract the first author and return only the author's name",
		//	"Extract the second author and return only the author's name",
		//	"Extract the year of publication and return only the year",
		//	"Extract a list of relevant keywords and return the keywords as a list, without any additional words or explanation",
		//	"Check if the document contains math formulas and return true or false",
		//},
	}

	tables := pb.NewTableServiceClient(conn)
	tableID, err := tables.CreateTable(ctx, table)
	if err != nil {
		log.Fatalf("did not create table: %v", err)
	}

	log.Printf("Table ID: %s", tableID.Id)

	column1, err := tables.AddColumnToTable(ctx, &pb.NewColumn{
		TableId:          tableID.Id,
		GenerationPrompt: "Extract the title",
	})
	if err != nil {
		log.Fatalf("did not create column: %v", err)
	}

	log.Printf("Column ID: %s", utils.Prettify(column1))

	//tester := test.NewTester(conn)
	//tester.TestDocumentRename()
	//tester.TestDocumentList()
	//tester.TestDocumentDelete()
	//tester.TestDocumentGet()
	//tester.TestWebpageIndex()
	//tester.TestStartThread()
	//tester.TestThreadMessages()
	//tester.TestGetThread()
	//tester.TestListThreadIDs()
	//tester.TestDeleteThread()
	//tester.TestDeleteThreadMessage()
	//tester.TestAccountCosts()
	//tester.TestGetCollection()
}
