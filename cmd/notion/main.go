package main

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm/bedrock"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"io"
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
		Prompt:       "List all Authors of the paper, separated by comma.",
		ModelOptions: &pb.ModelOptions{
			Model:       bedrock.ClaudeHaiku,
			Temperature: 1.0,
			MaxTokens:   256,
			TopP:        1.0,
		},
	})
	if err != nil {
		log.Fatalf("could not execute: %v", err)
	}

	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("could not receive: %v", err)
			break
		}

		log.Printf("Completed: %v", result.Document)
	}
}
