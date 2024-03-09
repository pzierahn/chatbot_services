package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	"io"
	"log"
	"net/http"
	"strings"
)

type Message struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

type RequestMeta struct {
	UserId string `json:"user_id,omitempty"`
}

type Request struct {
	Model         string      `json:"model,omitempty"`
	Messages      []Message   `json:"messages,omitempty"`
	System        string      `json:"system,omitempty"`
	MaxTokens     int         `json:"max_tokens,omitempty"`
	Metadata      RequestMeta `json:"metadata,omitempty"`
	StopSequences []string    `json:"stop_sequences,omitempty"`
	Stream        bool        `json:"stream,omitempty"`
	Temperature   float32     `json:"temperature,omitempty"`
	TopP          float32     `json:"top_p,omitempty"`
	TopK          int         `json:"top_k,omitempty"`
}

type Content struct {
	Text string `json:"text,omitempty"`
	Type string `json:"type,omitempty"`
}

type Response struct {
	Id         string    `json:"id,omitempty"`
	Model      string    `json:"model,omitempty"`
	Role       string    `json:"role,omitempty"`
	StopReason string    `json:"stop_reason,omitempty"`
	Type       string    `json:"type,omitempty"`
	Usage      Usage     `json:"usage,omitempty"`
	Content    []Content `json:"content,omitempty"`
}

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

func (client *Client) GenerateCompletion(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
	var messages []Message

	for _, msg := range req.Messages {
		var role string
		if msg.Type == llm.MessageTypeUser {
			role = RoleUser
		} else {
			role = RoleAssistant
		}

		if len(messages) > 0 && messages[len(messages)-1].Role == role {
			messages[len(messages)-1].Content += "\n" + msg.Text
		} else {
			messages = append(messages, Message{
				Role:    role,
				Content: msg.Text,
			})
		}
	}

	model := strings.TrimPrefix(req.Model, prefix)

	// Create a new request
	complete := &Request{
		Model:     model,
		Messages:  messages,
		System:    req.SystemPrompt,
		MaxTokens: req.MaxTokens,
		Metadata: RequestMeta{
			UserId: req.UserId,
		},
		Temperature: req.Temperature,
		TopP:        req.TopP,
	}

	out, err := json.MarshalIndent(complete, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Create a http request with custom headers
	httpClient := &http.Client{}

	endPoint := "https://api.anthropic.com/v1/messages"
	httpReq, err := http.NewRequest("POST", endPoint, bytes.NewReader(out))
	if err != nil {
		log.Fatal(err)
	}

	httpReq.Header.Add("x-api-key", client.apiKey)
	httpReq.Header.Add("anthropic-version", "2023-06-01")
	httpReq.Header.Add("content-type", "application/json")

	// Send the request
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		byt, err := io.ReadAll(resp.Body)
		if err == nil {
			log.Printf("Response: %v", string(byt))
		}

		return nil, fmt.Errorf("unexpected status code: %v", resp.Status)
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	usage := llm.ModelUsage{
		Model:        response.Model,
		UserId:       req.UserId,
		InputTokens:  response.Usage.InputTokens,
		OutputTokens: response.Usage.OutputTokens,
	}
	client.usage.Track(ctx, usage)

	return &llm.GenerateResponse{
		Text:  response.Content[0].Text,
		Usage: usage,
	}, nil
}
