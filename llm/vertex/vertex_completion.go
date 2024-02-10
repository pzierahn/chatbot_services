package vertex

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"github.com/pzierahn/chatbot_services/llm"
)

func (client *Client) GenerateCompletion(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
	modelName := req.Model
	if modelName == "" {
		modelName = "gemini-pro"
	}

	model := client.genaiClient.GenerativeModel(modelName)

	var parts []genai.Part
	for _, msg := range req.Messages {
		parts = append(parts, genai.Text(msg.Text))
	}

	gen, err := model.GenerateContent(ctx, parts...)
	if err != nil {
		return nil, err
	}

	if len(gen.Candidates) == 0 {
		return nil, nil
	}

	txt, ok := gen.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, nil
	}

	resp := &llm.GenerateResponse{
		Text: string(txt),
	}

	usage := llm.ModelUsage{
		UserId: req.UserId,
		Model:  modelName,
	}

	if gen.UsageMetadata != nil {
		usage.PromptTokens = int(gen.UsageMetadata.PromptTokenCount)
		usage.CompletionTokens = int(gen.UsageMetadata.CandidatesTokenCount)
	} else {
		for _, part := range parts {
			partText, ok := part.(genai.Text)
			if !ok {
				continue
			}
			usage.PromptTokens += len(string(partText))
		}
		usage.CompletionTokens = len(resp.Text)
	}

	client.trackUsage(ctx, usage)

	return resp, nil
}
