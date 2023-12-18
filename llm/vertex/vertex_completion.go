package vertex

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"github.com/pzierahn/brainboost/llm"
)

func (client *Client) GenerateCompletion(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
	modelName := req.Model
	if req.Model == "" {
		modelName = "gemini-pro"
	}

	model := client.genaiClient.GenerativeModel(modelName)

	var parts []genai.Part
	for _, part := range req.Documents {
		parts = append(parts, genai.Text(part))
	}

	parts = append(parts, genai.Text(req.Prompt))

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
		resp.InputTokens = int(gen.UsageMetadata.PromptTokenCount)
		resp.OutputTokens = int(gen.UsageMetadata.CandidatesTokenCount)
	}

	return resp, nil
}
