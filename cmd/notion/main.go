package main

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

	notion := pb.NewNotionClient(conn)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "User-Id", "j7jjxLD9rla2DrZoeUu3Tnft4812")

	stream, err := notion.ExecutePrompt(ctx, &pb.NotionPrompt{
		DatabaseID:   "8b9304529d664d2997834734345236f6",
		CollectionID: "59698763-c0ff-48c4-a69d-3d6ad62a7d50",
		Prompt:       "Give a list of evaluation metrics used in the training of a machine learning model",
	})

	for {
		result, err := stream.Recv()
		if err != nil {
			log.Fatalf("could not receive: %v", err)
		}

		log.Printf("Finish: %v", result)
	}
}
