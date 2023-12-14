package vertex

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
)

func (client *Client) Generate(ctx context.Context, prompt string, contexts ...string) (string, error) {
	model := client.genaiClient.GenerativeModel("gemini-pro")

	var parts []genai.Part
	for _, part := range contexts {
		parts = append(parts, genai.Text(part))
	}

	parts = append(parts, genai.Text(prompt))

	resp, err := model.GenerateContent(ctx, parts...)
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 {
		return "", nil
	}

	txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", nil
	}

	return string(txt), nil
}
