package vertex

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

func (client *Client) GenerateCompletion(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
	modelName, found := strings.CutPrefix(req.Model, modelPrefix)
	if !found || modelName == "" {
		modelName = "gemini-pro"
	}

	model := client.client.GenerativeModel(modelName)

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

	if gen.UsageMetadata != nil {
		client.usage.Track(ctx, llm.ModelUsage{
			UserId:           req.UserId,
			Model:            modelName,
			PromptTokens:     int(gen.UsageMetadata.PromptTokenCount),
			CompletionTokens: int(gen.UsageMetadata.CandidatesTokenCount),
		})
	}

	return resp, nil
}
