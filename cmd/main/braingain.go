package main

import (
	"braingain/pdf"
	"context"
	"flag"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
)

const baseDir = "/Users/patrick/patrick.zierahn@gmail.com - Google Drive/My Drive/KIT/2023-SS/DeSys/Lecture Slides"

var (
	//pdfFile     = baseDir + "/DeSys_11_Payment_Channel_Networks.pdf"
	pdfFile     = baseDir + "/DeSys_12_Smart_Contract_Platforms_Ethereum.pdf"
	first       = 8
	last        = 33
	temperature = 0.25
	prompt      = "What are smart contracts?"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.Parse()

	log.Printf("Filename: %s\n", pdfFile)

	pageCount, err := pdf.GetPageCount(context.Background(), pdfFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Page count: %d\n", pageCount)

	pages, err := pdf.ReadPages(context.Background(), pdfFile, pageCount)
	if err != nil {
		log.Fatal(err)
	}

	var messages []openai.ChatCompletionMessage

	if first > last {
		log.Fatal("First page must be less than or equal to last page")
	}

	for _, page := range pages[first-1 : last] {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: page,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	token := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(token)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo16K,
			Temperature: float32(temperature),
			Messages:    messages,
			N:           1,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	log.Println(resp.Choices[0].Message.Content)
	_ = os.WriteFile("output.txt", []byte(resp.Choices[0].Message.Content), 0644)
}
