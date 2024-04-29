package main

import (
	"context"
	"github.com/jomei/notionapi"
	"github.com/pzierahn/chatbot_services/llm/bedrock"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"strings"
)

// const databaseID = notionapi.DatabaseID("8b9304529d664d2997834734345236f6")
const databaseID = notionapi.DatabaseID("2705037dfb084e97b5ce578a497a5c34")

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

	token := os.Getenv("NOTION_API_KEY")
	client := notionapi.NewClient(notionapi.Token(token))

	ctx := context.Background()
	dbEntries, err := client.Database.Query(ctx, databaseID, &notionapi.DatabaseQueryRequest{
		Sorts: []notionapi.SortObject{
			{
				Property:  "ID",
				Direction: "ascending",
			},
		},
		PageSize: 999,
	})
	if err != nil {
		log.Fatal(err)
	}

	parallelPrompts := map[string]string{
		"subjects": "What is the number of subjects? Be concise and keep the answer short.",
	}

	var prompts []string
	var columns []string

	for column, prompt := range parallelPrompts {
		log.Printf("Create column %s", column)
		_, err = client.Database.Update(ctx, databaseID, &notionapi.DatabaseUpdateRequest{
			Properties: map[string]notionapi.PropertyConfig{
				column: notionapi.RichTextPropertyConfig{
					Type: notionapi.PropertyConfigTypeRichText,
				},
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		prompts = append(prompts, prompt)
		columns = append(columns, column)
	}

	log.Printf("Collect Notion page-ids")

	// Document Name --> Notion PageID
	pageIDs := make(map[string]notionapi.PageID)
	for _, result := range dbEntries.Results {
		props := result.Properties

		rich, ok := props["ID"].(*notionapi.TitleProperty)
		if !ok {
			continue
		}

		if len(rich.Title) <= 0 {
			continue
		}

		pageIDs[rich.Title[0].PlainText] = notionapi.PageID(result.ID)
	}

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

	log.Printf("List collection documents")
	documentList, err := documentService.List(ctxx, &pb.DocumentFilter{
		CollectionId: "59698763-c0ff-48c4-a69d-3d6ad62a7d50",
	})
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

	resp, err := chatService.BatchChat(ctxx, &pb.BatchRequest{
		DocumentIds: docIDs,
		Prompts:     prompts,
		ModelOptions: &pb.ModelOptions{
			Model:       bedrock.ClaudeHaiku,
			Temperature: 1,
			MaxTokens:   256,
			TopP:        1,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	//
	// Store the results on notion
	//

	log.Printf("Store results in Notion: %d", len(resp.Items))

	for _, completion := range resp.Items {
		docID := resp.DocumentIds[completion.DocumentId]
		docName := idsName[docID]
		pageID := pageIDs[strings.TrimSuffix(docName, ".pdf")]

		column := columns[completion.Prompt]

		log.Println(docName, column)

		_, err = client.Page.Update(ctx, pageID, &notionapi.PageUpdateRequest{
			Properties: map[string]notionapi.Property{
				column: notionapi.RichTextProperty{
					Type: notionapi.PropertyTypeRichText,
					RichText: []notionapi.RichText{
						{
							Text: &notionapi.Text{
								Content: completion.Completion,
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Println(err)
		}
	}
}
