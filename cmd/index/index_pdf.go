package main

import (
	"braingain/database"
	"braingain/pdf"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/qdrant/go-client/qdrant"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	baseDir    = "/Users/patrick/patrick.zierahn@gmail.com - Google Drive/My Drive/KIT/2023-SS/DeSys/Lecture Slides"
	collection = "DeSys"
)

type queryResult struct {
	Id       string
	Score    float32
	Filename string
	Page     int
	Content  string
}

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

		// Create hash for file
		hash := sha256.New()
		hash.Write([]byte(file.Name()))

		log.Printf("--> Pages: %d\n", len(pages))

		for inx, page := range pages {
			log.Printf("--> %d/%d\n", inx+1, len(pages))

			page = strings.TrimSpace(page)
			if len(page) == 0 {
				continue
			}

			resp, err := ai.CreateEmbeddings(
				context.Background(),
				openai.EmbeddingRequestStrings{
					Model: openai.AdaEmbeddingV2,
					Input: []string{page},
				},
			)

			if err != nil {
				log.Fatalf("could not create embeddings: %v", err)
			}

			err = conn.Upsert(ctx, database.Payload{
				Uuid:       uuid.NewString(),
				Collection: collection,
				Data:       resp.Data[0].Embedding,
				Metadata: map[string]*pb.Value{
					"hash": {
						Kind: &pb.Value_StringValue{
							StringValue: fmt.Sprintf("%x", hash.Sum(nil)),
						},
					},
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

	//search := "Explain the Practical Byzantine Fault Tolerance (PBFT) algorithm in detail"
	search := "What is the difference between the Practical Byzantine Fault Tolerance (PBFT) algorithm and the DAG-Rider algorithm?"

	token := os.Getenv("OPENAI_API_KEY")
	ai := openai.NewClient(token)
	resp, err := ai.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequestStrings{
			Model: openai.AdaEmbeddingV2,
			Input: []string{search},
		},
	)

	if err != nil {
		log.Fatalf("could not create embeddings: %v", err)
	}

	// Search for a page
	embedding := resp.Data[0].Embedding

	searchResponse, err := conn.SearchEmbedding(ctx, collection, embedding)
	if err != nil {
		log.Fatalf("could not search points: %v", err)
	}

	pages := make(map[string][]int)
	queryResults := make([]queryResult, 0)

	log.Printf("Search: %v\n", searchResponse.Time)
	for _, hit := range searchResponse.Result {
		if hit.Payload["filename"].GetStringValue() == "DeSys_07_Consensus_and_Variants_v1.pdf" {
			continue
		}

		filename := hit.Payload["filename"].GetStringValue()
		page := int(hit.Payload["page"].GetIntegerValue()) + 1

		if _, ok := pages[filename]; !ok {
			pages[filename] = make([]int, 0)
		}

		pages[filename] = append(pages[filename], page)

		//log.Printf("Hit:      %v\n", hit.Id)
		//log.Printf("Score:    %v\n", hit.Score)
		//log.Printf("Filename: %v\n", filename)
		//log.Printf("Page:     %v\n", page)
		//log.Printf("-------------\n")

		queryResults = append(queryResults, queryResult{
			Id:       hit.Id.GetUuid(),
			Score:    hit.Score,
			Filename: filename,
			Page:     page,
			Content:  hit.Payload["content"].GetStringValue(),
		})
	}

	sort.SliceStable(queryResults, func(i, j int) bool {
		return queryResults[i].Page < queryResults[j].Page
	})
	sort.SliceStable(queryResults, func(i, j int) bool {
		return queryResults[i].Filename < queryResults[j].Filename
	})

	byt, _ := json.MarshalIndent(queryResults, "", "  ")
	log.Printf("Query results: %s\n", byt)

	err = os.WriteFile("query_results.json", byt, 0644)
	if err != nil {
		log.Fatalf("could not write query results: %v", err)
	}

	for filename, pages := range pages {
		sort.Ints(pages)
		log.Printf("%v --> %v\n", filename, pages)
	}

	var messages []openai.ChatCompletionMessage

	for _, result := range queryResults {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: result.Content,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: search,
	})

	chat, err := ai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			//Model:       openai.GPT3Dot5Turbo16K,
			Model:       openai.GPT4,
			Temperature: float32(0),
			Messages:    messages,
			N:           1,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	usage, _ := json.MarshalIndent(chat.Usage, "", "  ")
	log.Printf("Usage: %s\n", usage)

	log.Println(chat.Choices[0].Message.Content)
	_ = os.WriteFile("output.txt", []byte(chat.Choices[0].Message.Content), 0644)
}
