package main

import (
	"braingain/pdf"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
)

const baseDir = "/Users/patrick/patrick.zierahn@gmail.com - Google Drive/My Drive/KIT/2023-SS/DeSys/Lecture Slides"

var (
	//pdfFile = baseDir + "/DeSys_01_Intro.pdf"
	//pdfFile = baseDir + "/DeSys_09_Decentralized_Messaging_Matrix.pdf"
	//pdfFile = baseDir + "/DeSys_07_Consensus_and_Variants_v2.pdf"
	//pdfFile = baseDir + "/DeSys_11_Payment_Channel_Networks.pdf"
	//pdfFile = baseDir + "/DeSys_12_Smart_Contract_Platforms_Ethereum.pdf"
	//pdfFile     = baseDir + "/DeSys_13_Decentralized_File_Storage_IPFS.pdf"
	//pdfFile     = "/Users/patrick/patrick.zierahn@gmail.com - Google Drive/My Drive/KIT/2023-SS/DeSys/Further Readings/176429260X.pdf"
	//pdfFile     = "/Users/patrick/patrick.zierahn@gmail.com - Google Drive/My Drive/KIT/2023-SS/DeSys/Further Readings/Harvest_yield_and_scalable_tolerant_systems.pdf"
	//pdfFile = "/Users/patrick/patrick.zierahn@gmail.com - Google Drive/My Drive/KIT/2023-SS/DeSys/Further Readings/1-s2.0-089054018790054X-main.pdf"
	pdfFile = "/Users/patrick/patrick.zierahn@gmail.com - Google Drive/My Drive/KIT/2023-SS/DeSys/Further Readings/2102.08325.pdf"
	first   = 1
	last    = 8
	//first       = 4
	//last        = 20
	temperature = 0.0
	//model       = openai.GPT3Dot5Turbo16K
	model = openai.GPT3Dot5Turbo16K
	//prompt = "Explain the DAG-Rider algorithm"
	//prompt = "How does leader election work in DAG-Rider?"
	prompt = "How is the global perfect coin in DAG-Rider determined?"
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
