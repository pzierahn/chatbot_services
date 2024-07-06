package vertex

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

func (client *Client) Completion(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	modelName, _ := strings.CutPrefix(req.Model, modelPrefix)

	outputTokens := int32(req.MaxTokens)

	model := client.client.GenerativeModel(modelName)
	model.TopP = &req.TopP
	model.Temperature = &req.Temperature
	model.MaxOutputTokens = &outputTokens
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(req.SystemPrompt)},
	}
	model.Tools = client.getTools()

	chat := model.StartChat()

	var parts []genai.Part
	for _, msg := range req.Messages {
		parts = append(parts, genai.Text(msg.Content))
	}

	gen, err := chat.SendMessage(ctx, parts...)
	if err != nil {
		return nil, err
	}

	if len(gen.Candidates) == 0 || len(gen.Candidates[0].Content.Parts) == 0 {
		return nil, nil
	}

	usage := llm.ModelUsage{
		UserId: req.UserId,
		Model:  modelName,
	}

	if gen.UsageMetadata != nil {
		usage.InputTokens = int(gen.UsageMetadata.PromptTokenCount)
		usage.OutputTokens = int(gen.UsageMetadata.CandidatesTokenCount)
	}

	if fun, ok := gen.Candidates[0].Content.Parts[0].(genai.FunctionCall); ok {
		parts = append(parts, fun)

		result, err := client.callTool(ctx, fun.Name, fun.Args)
		if err != nil {
			return nil, err
		}

		functionResults := genai.FunctionResponse{
			Name: fun.Name,
			Response: map[string]any{
				"content": result,
			},
		}
		parts = append(parts, functionResults)

		gen, err = chat.SendMessage(ctx, functionResults)
		if err != nil {
			return nil, err
		}

		if gen.UsageMetadata != nil {
			usage.InputTokens += int(gen.UsageMetadata.PromptTokenCount)
			usage.OutputTokens += int(gen.UsageMetadata.CandidatesTokenCount)
		}
	}

	txt, ok := gen.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, nil
	}

	return &llm.CompletionResponse{
		Message: &llm.Message{
			Role:    llm.MessageTypeUser,
			Content: string(txt),
		},
		Usage: usage,
	}, nil
}
