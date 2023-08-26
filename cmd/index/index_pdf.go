package main

import (
	"braingain/database"
	"braingain/pdf"
	"context"
	"github.com/google/uuid"
	pb "github.com/qdrant/go-client/qdrant"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"strings"
)

const (
	baseDir    = "/Users/patrick/patrick.zierahn@gmail.com - Google Drive/My Drive/KIT/2023-SS/DeSys/Lecture Slides"
	collection = "DeSys"
)

func indexFiles(conn *database.Client) {
	ctx := context.Background()

	_ = conn.DeleteCollection(ctx, collection)
	err := conn.CreateCollection(ctx, collection, 1536, pb.Distance_Cosine)
	if err != nil {
		log.Fatalf("could not create collection: %v", err)
	}

	// Read PDF files in baseDir
	files, err := os.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}

	token := os.Getenv("OPENAI_API_KEY")
	ai := openai.NewClient(token)

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".pdf") {
			continue
		}

		log.Printf("Filename: %s\n", file.Name())

		pdfFile := baseDir + "/" + file.Name()
		pages, err := pdf.ReadPages(ctx, pdfFile)
		if err != nil {
			log.Fatalf("could not read pages: %v", err)
		}

		for inx, page := range pages {
			resp, err := ai.CreateEmbeddings(
				context.Background(),
				openai.EmbeddingRequestStrings{
					Model: openai.AdaEmbeddingV2,
					Input: []string{
						page,
					},
				},
			)

			//data := make([]float32, 1536)
			//data[0] = rand.Float32()

			err = conn.Upsert(ctx, database.Payload{
				Uuid:       uuid.NewString(),
				Collection: collection,
				Data:       resp.Data[0].Embedding,
				Metadata: map[string]*pb.Value{
					"filename": {
						Kind: &pb.Value_StringValue{
							StringValue: file.Name(),
						},
					},
					"page": {
						Kind: &pb.Value_IntegerValue{
							IntegerValue: int64(inx),
						},
					},
					"content": {
						Kind: &pb.Value_StringValue{
							StringValue: page,
						},
					},
				},
			})

			if err != nil {
				log.Fatalf("could not upsert points: %v", err)
			}
		}
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	conn, err := database.Connect("localhost:6334")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() { _ = conn.Close() }()

	// indexFiles(conn)

	ctx := context.Background()
	count, err := conn.Count(ctx, collection)
	if err != nil {
		log.Fatalf("could not count points: %v", err)
	}

	log.Printf("Count: %v\n", count.Result.Count)

	// Search for a page
	embedding := make([]float32, 1536)
	embedding[0] = 0.5

	searchResponse, err := conn.SearchEmbedding(ctx, collection, embedding)
	if err != nil {
		log.Fatalf("could not search points: %v", err)
	}

	log.Printf("Search: %v\n", searchResponse.Time)
	for _, hit := range searchResponse.Result {
		log.Printf("Hit:      %v\n", hit.Id)
		log.Printf("Filename: %v\n", hit.Payload["filename"])
		log.Printf("Page:     %v\n", hit.Payload["page"])
	}
}
