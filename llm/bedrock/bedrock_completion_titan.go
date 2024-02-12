package bedrock

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

type TitanTextGenerationConfig struct {
	MaxTokenCount int      `json:"maxTokenCount"`
	StopSequence  []string `json:"stopSequences"`
	Temperature   float64  `json:"temperature"`
	TopP          float64  `json:"topP"`
}

type TitanRequest struct {
	Prompt string                    `json:"inputText"`
	Config TitanTextGenerationConfig `json:"textGenerationConfig"`
}

type TitanResponse struct {
	InputTextTokenCount int `json:"inputTextTokenCount"`
	Results             []struct {
		TokenCount       int    `json:"tokenCount"`
		OutputText       string `json:"outputText"`
		CompletionReason string `json:"completionReason"`
	} `json:"results"`
}

func (client *Client) generateCompletionTitan(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
	var parts []string
	for _, msg := range req.Messages {
		if msg.Type == llm.MessageTypeUser {
			parts = append(parts, "Human: "+msg.Text)
		} else {
			parts = append(parts, "Bot: "+msg.Text)
		}
	}

	request := TitanRequest{
		Prompt: strings.Join(parts, "\n\n"),
		Config: TitanTextGenerationConfig{
			MaxTokenCount: 8192,
			StopSequence:  []string{},
			TopP:          1,
		},
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	result, err := client.bedrock.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(req.Model),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        jsonBody,
	})
	if err != nil {
		return nil, err
	}

	var response TitanResponse
	err = json.Unmarshal(result.Body, &response)
	if err != nil {
		return nil, err
	}

	return &llm.GenerateResponse{
		Text: response.Results[0].OutputText,
	}, nil
}
