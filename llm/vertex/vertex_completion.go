package vertex

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

func (client *Client) GenerateCompletion(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
	modelName, _ := strings.CutPrefix(req.Model, modelPrefix)

	outputTokens := int32(req.MaxTokens)

	model := client.client.GenerativeModel(modelName)
	model.TopP = &req.TopP
	model.Temperature = &req.Temperature
	model.MaxOutputTokens = &outputTokens

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
		resp.Usage = llm.ModelUsage{
			UserId:       req.UserId,
			Model:        modelName,
			InputTokens:  int(gen.UsageMetadata.PromptTokenCount),
			OutputTokens: int(gen.UsageMetadata.CandidatesTokenCount),
		}

		client.usage.Track(ctx, resp.Usage)
	}

	return resp, nil
}
