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
	pis      = "/Users/patrick/patrick.zierahn@gmail.com - Google Drive/My Drive/KIT/2023-SS/Praktikum Ingenieursmäßige Software-Entwicklung/"
)

var (
	//pdfFile = baseDir + "/DeSys_09_Decentralized_Messaging_Matrix.pdf"
	//pdfFile = baseDir + "/DeSys_08_Consistency.pdf"
	//pdfFile     = readings + "/IPTPS2002.pdf"
	pdfFile     = pis + "/Evaluation_Methods_and_Replicability_of_Software_Architecture_Research_Objects.pdf"
	first       = 1
	last        = 10
	temperature = 0.0
	model       = openai.GPT3Dot5Turbo16K
	//model = openai.GPT4
	prompt = "Build bullet points for a single presentation slide."
	//prompt = "Explain the paper"
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
