package main

import (
	"context"
	"encoding/csv"
	"github.com/jomei/notionapi"
	"github.com/pzierahn/chatbot_services/llm/bedrock"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"strconv"
	"strings"
)

const databaseID = notionapi.DatabaseID("8b9304529d664d2997834734345236f6")

func readCSV() [][]string {
	// os.Open() opens specific file in
	// read-only mode and this return
	// a pointer of type os.File
	file, err := os.Open("/Users/patrick/Downloads/Literature Research - Taxonomy.csv")

	// Checks for the error
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	// Closes the file
	defer func() {
		_ = file.Close()
	}()

	// The csv.NewReader() function is called in
	// which the object os.File passed as its parameter
	// and this creates a new csv.Reader that reads
	// from the file
	reader := csv.NewReader(file)

	// ReadAll reads all the records from the CSV file
	// and Returns them as slice of slices of string
	// and an error if any
	records, err := reader.ReadAll()

	// Checks for the error
	if err != nil {
		log.Fatal("Error reading records")
	}

	return records
}

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

func importCSV(client notionapi.Client) {
	records := readCSV()
	for _, record := range records[1:] {
		var (
			name     = record[0]
			title    = record[1]
			year     = record[2]
			approach = record[3]
			method   = record[4]
			metrics  = record[5]
			results  = record[6]
			subjects = record[7]
			roi      = record[8]
			neuronal = record[9]
			//implement  = record[10]
			dataset = record[11]
		)

		log.Println(name)

		if neuronal == "" {
			neuronal = "TODO"
		}

		if roi == "" {
			roi = "TODO"
		}

		// Parse year to a float
		yearFloat, _ := strconv.ParseFloat(year, 32)

		ctx := context.Background()
		_, err := client.Page.Create(ctx, &notionapi.PageCreateRequest{
			Parent: notionapi.Parent{
				DatabaseID: databaseID,
			},
			Properties: map[string]notionapi.Property{
				"Name": notionapi.TitleProperty{
					Type: notionapi.PropertyTypeTitle,
					Title: []notionapi.RichText{
						{
							Text: &notionapi.Text{
								Content: name,
							},
						},
					},
				},
				"Title": notionapi.RichTextProperty{
					Type: notionapi.PropertyTypeRichText,
					RichText: []notionapi.RichText{
						{
							Text: &notionapi.Text{
								Content: title,
							},
						},
					},
				},
				"Approach": notionapi.MultiSelectProperty{
					Type: notionapi.PropertyTypeMultiSelect,
					MultiSelect: []notionapi.Option{
						{
							Name: approach,
						},
					},
				},
				"Year": notionapi.NumberProperty{
					Type:   notionapi.PropertyTypeNumber,
					Number: yearFloat,
				},
				"Algorithms and Method": notionapi.RichTextProperty{
					Type: notionapi.PropertyTypeRichText,
					RichText: []notionapi.RichText{
						{
							Text: &notionapi.Text{
								Content: method,
							},
						},
					},
				},
				"Evaluation Metrics": notionapi.RichTextProperty{
					Type: notionapi.PropertyTypeRichText,
					RichText: []notionapi.RichText{
						{
							Text: &notionapi.Text{
								Content: metrics,
							},
						},
					},
				},
				"Evaluation Results": notionapi.RichTextProperty{
					Type: notionapi.PropertyTypeRichText,
					RichText: []notionapi.RichText{
						{
							Text: &notionapi.Text{
								Content: results,
							},
						},
					},
				},
				"Number of Subjects": notionapi.RichTextProperty{
					Type: notionapi.PropertyTypeRichText,
					RichText: []notionapi.RichText{
						{
							Text: &notionapi.Text{
								Content: subjects,
							},
						},
					},
				},
				"ROI": notionapi.SelectProperty{
					Type: notionapi.PropertyTypeSelect,
					Select: notionapi.Option{
						Name: roi,
					},
				},
				"Neuronal Networks": notionapi.SelectProperty{
					Type: notionapi.PropertyTypeSelect,
					Select: notionapi.Option{
						Name: neuronal,
					},
				},
				"Dataset": notionapi.RichTextProperty{
					Type: notionapi.PropertyTypeRichText,
					RichText: []notionapi.RichText{
						{
							Text: &notionapi.Text{
								Content: dataset,
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		//break
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	token := os.Getenv("NOTION_API_KEY")
	client := notionapi.NewClient(notionapi.Token(token))

	ctx := context.Background()
	dbEntries, err := client.Database.Query(ctx, databaseID, &notionapi.DatabaseQueryRequest{
		Sorts: []notionapi.SortObject{
			{
				Property:  "Title",
				Direction: "ascending",
			},
		},
		PageSize: 999,
	})
	if err != nil {
		log.Fatal(err)
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
		Prompts: []string{
			//"List the used Datasets, use only the abbreviation. Be concise and keep the answer short",
			//"Create list of used algorithms. Be concise and keep the answer short",
			//"Does the paper mention a GitHub repository? Just answer Yes or No",
			//"Extract the MAE, RMAE, RMSE, MSE and Pearson Correlation Coefficient. " +
			//	"Be concise and keep the answer short",
			//"Classify how the ROI is selected: Automatically, Manually, None",
			//"Extract all GitHubs urls. If no exist return \"-\"",
			"Is the proposed respiratory extraction method based on PPG or motion?",
		},
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

		log.Println(docID, docName, pageID)

		_, err := client.Page.Update(ctx, pageID, &notionapi.PageUpdateRequest{
			Properties: map[string]notionapi.Property{
				"Motion/rPPG": notionapi.RichTextProperty{
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
