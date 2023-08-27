package main

import (
	"braingain/database"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"sort"
)

const (
	collection = "DeSys"
)

type queryResult struct {
	Id       string
	Score    float32
	Filename string
	Page     int
	Content  string
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	conn, err := database.Connect("localhost:6334")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() { _ = conn.Close() }()

	ctx := context.Background()

	//search := "Explain the Practical Byzantine Fault Tolerance (PBFT) algorithm in detail"
	//search := "What is the difference between the Practical Byzantine Fault Tolerance (PBFT) algorithm and the DAG-Rider algorithm?"
	//search := "How does the DAG-Rider algorithm work?"
	//search := "Why does the DAG-Rider algorithm needs waves?"
	//search := "Explain the DAG-Rider algorithm in detail"
	//search := "Which properties of an Operation-Based CRDT have to be shown?"
	//search := "How does a Sybil Attack with Distributed Secret Sharing work?"
	//search := "How is consensus archived in total order broadcasts?"
	//search := "What is the differance between consistency and consensus?"
	//search := "How does the TEE-Rider algorithm work in detail?"
	//search := "What is the TEE-Rider algorithm and how does it work in detail?"
	//search := "What is a Byzantine Atomic Broadcast?"
	//search := "What is the difference between a Byzantine Atomic Broadcast and a Byzantine Reliable Broadcast?"
	//search := "How can total order be archived with a Byzantine Atomic Broadcast?"
	//search := "What is the definition of Byzantine Broadcast Channel? What is the difference between RB-Agreement and RB-Validity?"
	//search := "Explain the TEE-based Reliable Broadcast setting"
	//search := "Explain the TEE-based Reliable Broadcast in detail"
	//search := "Explain the TEE-Rider algorithm in detail"
	//search := "What are the failure models?"
	//search := "How can correct processors verify that a number was generated only by a Unique Sequential Identifier Generator (USIG)?"
	//search := "How is a message signed with USIG?"
	//search := "What is the difference between DAG-Rider and TEE-Rider algorithms?"
	//search := "What is MinBFT?"
	//search := "In how far does the TEE-based Reliable Broadcast differ from other reliable broadcasts?"
	//search := "How is partial synchrony, synchrony and asynchrony defined?"
	//search := "What is the difference between partial synchrony and asynchrony defined?"
	//search := "What is the advantage of knowing that \"Δ is arbitrary but fix and unknown\" in partial synchrony?"
	//search := "How does the TEE-Rider algorithm uses four waves?"
	//search := "What is the advantage of partial synchrony against Asynchrony?"
	//search := "Why can partial synchrony guarantee safety and liveness?"
	//search := "Explain the flp impossibility in detail"
	//search := "Explain the impossibility of Distributed Consensus with One Faulty Process"
	//search := "How can a Total Order Broadcast be archived?"
	//search := "How can consensus be derived in TEE-Rider?"
	//search := "How is a Total Order in TEE-Rider archived?"
	//search := "What makes a distributed system a decentralized system? → The Meaning of “Decentralization”"
	//search := "What are the reasons for decentralization?"
	//search := "What is the CAP Theorem? Why can only two of the three properties be archived?"
	//search := "What is the CAP Theorem? Why can only two of the three properties be archived? Give an examples"
	//search := "What is the differance between Consistency and Partition Tolerance in CAP?"
	//search := "Why is there no terminating algorithm for consensus in the asynchronous model?"
	//search := "Why it is impossible for processes in an asynchronous system to unanimously agree on a consensus value if even a single process could fail?"
	//search := "Explain the FLP Impossibility in detail"
	//search := "What is the difference between synchrony, partial synchrony and asynchrony? Is the delta in partial synchrony approximated?"
	//search := "What are the implications if the delta in partial synchrony approximated wrongly?"
	//search := "Which properties need to be shown to proof eventual consistency and strong eventual consistency?"
	//search := "How does leader election work in PBFT?"
	//search := "Why are 2f+1 needed to deal with f faulty processes?"
	//search := "How can f be determined in DAG-Rider?"
	//search := "What is deterministic threshold signatures? How do they work?"
	//search := "How does deterministic threshold signing work in DAG-Rider?"
	//search := "How does deterministic threshold signing work?"
	//search := "How does Distributed Key Generation (DKG) for RSA work?"
	//search := "How is safety and liveliness archived in DAG-Rider?"
	//search := "How does leader election in DAG-Rider work?"
	search := "Why are 4 rounds a wave in DAG-Rider?"

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
		filename := hit.Payload["filename"].GetStringValue()
		page := int(hit.Payload["page"].GetIntegerValue()) + 1

		if _, ok := pages[filename]; !ok {
			pages[filename] = make([]int, 0)
		}

		pages[filename] = append(pages[filename], page)

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
	//log.Printf("Query results: %s\n", byt)

	err = os.WriteFile("query_results.json", byt, 0644)
	if err != nil {
		log.Fatalf("could not write query results: %v", err)
	}

	// sort pages keys
	var files []string
	for filename := range pages {
		files = append(files, filename)
	}
	sort.Strings(files)

	for _, filename := range files {
		filepages := pages[filename]
		sort.Ints(filepages)
		log.Printf("%v --> %v\n", filename, filepages)
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
			//Model: openai.GPT3Dot5Turbo16K,
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
