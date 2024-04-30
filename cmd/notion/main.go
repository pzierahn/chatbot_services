package main

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm/bedrock"
	notion2 "github.com/pzierahn/chatbot_services/notion"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
	"sync"
)

// const databaseID = "8b9304529d664d2997834734345236f6"
const databaseID = "2705037dfb084e97b5ce578a497a5c34"

func findDocumentIDs(list *pb.DocumentList) (map[string]string, map[string]string) {
	// Map document names to document IDs
	nameIds := make(map[string]string)
	idsName := make(map[string]string)

	for docID, document := range list.Items {
		file := document.GetFile()
		nameIds[file.Filename] = docID
		idsName[docID] = file.Filename
	}

	return nameIds, idsName
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	notion, err := notion2.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	column := "eval"
	prompt := "Give a list of evaluation metrics mentioned in the document."

	err = notion.AddColumn(ctx, databaseID, column)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Collect Notion page-ids")

	conn, err := grpc.Dial(
		"localhost:8869",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() { _ = conn.Close() }()

	documentService := pb.NewDocumentServiceClient(conn)
	chatService := pb.NewChatServiceClient(conn)

	ctxx := metadata.AppendToOutgoingContext(ctx, "User-Id", "j7jjxLD9rla2DrZoeUu3Tnft4812")
	ctxx, cnl := context.WithCancel(ctxx)
	defer cnl()

	log.Printf("List collection documents")
	documentList, err := documentService.List(ctxx, &pb.DocumentFilter{
		CollectionId: "59698763-c0ff-48c4-a69d-3d6ad62a7d50",
	})
	if err != nil {
		log.Fatal(err)
	}

	pageIDs, err := notion.ListDocumentIDs(ctx, databaseID)
	if err != nil {
		log.Fatal(err)
	}

	nameIds, idsName := findDocumentIDs(documentList)
	var docIDs []string
	for docName := range pageIDs {
		docID, ok := nameIds[docName+".pdf"]
		if !ok {
			continue
		}

		docIDs = append(docIDs, docID)
	}

	log.Printf("Batch query execution on %d documents", len(docIDs))

	funcs := make([]func() (string, string, error), 0)

	for _, docID := range docIDs {
		funcs = append(funcs, func() (string, string, error) {
			resp, errx := chatService.Completion(ctxx, &pb.CompletionRequest{
				DocumentId: docID,
				Prompt:     prompt,
				ModelOptions: &pb.ModelOptions{
					Model:       bedrock.ClaudeHaiku,
					Temperature: 1,
					MaxTokens:   256,
					TopP:        1,
				},
			})

			return docID, resp.Completion, errx
		})
	}

	var wg sync.WaitGroup
	wg.Add(len(funcs))

	syn := make(chan struct{}, 10)

	for idx := range len(funcs) {
		go func(idx int) {
			defer wg.Done()

			syn <- struct{}{}
			defer func() { <-syn }()

			docID, completion, err := funcs[idx]()
			if err != nil {
				return
			}

			docName := idsName[docID]
			docName = strings.TrimSuffix(docName, ".pdf")

			log.Printf("Insert %s", docName)
			err = notion.UpdatePage(ctx, pageIDs[docName], column, completion)
			if err != nil {
				log.Fatal(err)
			}
		}(idx)
	}

	wg.Wait()
}
