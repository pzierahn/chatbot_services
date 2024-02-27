package bedrock

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

type ClaudeRequest struct {
	Prompt            string   `json:"prompt"`
	MaxTokensToSample int      `json:"max_tokens_to_sample"`
	Temperature       float64  `json:"temperature,omitempty"`
	TopP              float64  `json:"top_p,omitempty"`
	TopK              int      `json:"top_k,omitempty"`
	StopSequences     []string `json:"stop_sequences,omitempty"`
}

type ClaudeResponse struct {
	Completion string `json:"completion"`
}

func (client *Client) generateCompletionAnthropic(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
	prompt := ""

	for _, msg := range req.Messages {
		switch msg.Type {
		case llm.MessageTypeSystem:
			prompt += "\n\nSystem: " + msg.Text
		case llm.MessageTypeUser:
			prompt += "\n\nHuman: " + msg.Text
		case llm.MessageTypeBot:
			prompt += "\n\nAssistant: " + msg.Text
		}
	}

	prompt += "\n\nAssistant: "
	prompt = strings.TrimSpace(prompt)

	request := ClaudeRequest{
		Prompt:            prompt,
		MaxTokensToSample: 1024,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	result, err := client.bedrock.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(req.Model),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return nil, err
	}

	var response ClaudeResponse
	err = json.Unmarshal(result.Body, &response)
	if err != nil {
		return nil, err
	}

	return &llm.GenerateResponse{
		Text: strings.TrimSpace(response.Completion),
	}, nil
}
