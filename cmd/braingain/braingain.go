package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pzierahn/braingain/pdf"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"strings"
)

const (
	desysDir = "/Users/patrick/patrick.zierahn@gmail.com - Google Drive/My Drive/KIT/2023-SS/DeSys"
	baseDir  = desysDir + "/Lecture Slides"
	readings = desysDir + "/Further Readings"
)

var (
	//pdfFile = baseDir + "/DeSys_10_Distributed_Ledgers_Blockchains_Bitcoin.pdf"
	pdfFile     = readings + "/sigma.pdf"
	first       = 1
	last        = 7
	temperature = 1.0
	model       = openai.GPT3Dot5Turbo16K
	//model = openai.GPT4
	//prompt = "Read section 4 of the paper. How do the additions compared to the strawman protocol from section 3 prevent the scenario from the previous question from happening? Why does the informed backoff not run into a similar problem as the starwman protocol?"
	//prompt = "Why can the authors model a crash fault with loss of internal state?"
	//prompt = "Explain the algorithm"
	prompt = "Explain the terms 'informed backoff' and 'lease'"
	//prompt = "Explain what happens if two processes requesting the same resource concurrently"
	//prompt = "Explain how failure recovery works"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.Parse()

	log.Printf("Filename: %s\n", pdfFile)

	pages, err := pdf.ReadPages(context.Background(), pdfFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Page count: %d\n", len(pages))

	if first > last {
		log.Fatal("First page must be less than or equal to last page")
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: strings.Join(pages[first-1:last], "\n"),
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	token := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(token)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       model,
			Temperature: float32(temperature),
			Messages:    messages,
			N:           1,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	usage, _ := json.MarshalIndent(resp.Usage, "", "  ")
	log.Printf("Usage: %s\n", usage)

	log.Println(resp.Choices[0].Message.Content)
	_ = os.WriteFile("output.txt", []byte(resp.Choices[0].Message.Content), 0644)
}
