package main

import (
	"context"
	"encoding/csv"
	"github.com/jomei/notionapi"
	pb "github.com/pzierahn/chatbot_services/proto"
	"log"
	"os"
	"strconv"
)

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

func findDocumentIDs(list *pb.DocumentList) map[string]string {
	// Map document names to document IDs
	matches := make(map[string]string)

	for docID, document := range list.Items {
		file := document.GetFile()
		matches[file.Filename] = docID
	}

	return matches
}

func connect() {
	//conn, err := grpc.Dial(
	//	"localhost:8869",
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//)
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer func() { _ = conn.Close() }()
	//
	//documentService := pb.NewDocumentServiceClient(conn)
	//chatService := pb.NewChatServiceClient(conn)
	//
	//ctx := context.Background()
	//ctx = metadata.AppendToOutgoingContext(ctx, "User-Id", "j7jjxLD9rla2DrZoeUu3Tnft4812")
	//
	//documentList, err := documentService.List(ctx, &pb.DocumentFilter{
	//	CollectionId: "59698763-c0ff-48c4-a69d-3d6ad62a7d50",
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	////log.Print(documentList)
	//
	//docIDs := findDocumentIDs(documentList)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()
	databaseID := notionapi.DatabaseID("")
	token := os.Getenv("NOTION_API_KEY")
	client := notionapi.NewClient(notionapi.Token(token))

	//db, err := client.Database.Update(ctx, databaseID, &notionapi.DatabaseUpdateRequest{
	//	Properties: notionapi.PropertyConfigs{},
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println()

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
