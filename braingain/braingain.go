package braingain

import (
	"context"
	"github.com/pzierahn/braingain/database_pg"
	"github.com/sashabaranov/go-openai"
	"sort"
)

const (
	collection = "DeSys"
)

type Costs struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	PromptCosts      float32
	CompletionCosts  float32
	TotalCosts       float32
}

type ChatCompletion struct {
	Completion string
	Sources    []database_pg.ScorePoints
	Costs      Costs
}

type Chat struct {
	db    *database_pg.Client
	gpt   *openai.Client
	Model string
}

func NewChat(db *database_pg.Client, gpt *openai.Client) *Chat {
	return &Chat{
		db:    db,
		gpt:   gpt,
		Model: openai.GPT3Dot5Turbo16K,
	}
}

func (chat Chat) createEmbedding(ctx context.Context, prompt string) ([]float32, error) {
	resp, err := chat.gpt.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequestStrings{
			Model: openai.AdaEmbeddingV2,
			Input: []string{prompt},
		},
	)

	if err != nil {
		return nil, err
	}

	return resp.Data[0].Embedding, nil
}

func (chat Chat) calculateCosts(usage openai.Usage) Costs {
	costs := Costs{
		PromptTokens:     usage.PromptTokens,
		CompletionTokens: usage.CompletionTokens,
		TotalTokens:      usage.TotalTokens,
	}

	inputTokens := float32(usage.PromptTokens) / float32(1000)
	outputTokens := float32(usage.CompletionTokens) / float32(1000)

	switch chat.Model {
	case openai.GPT3Dot5Turbo:
		costs.PromptCosts = inputTokens * 0.0015
		costs.CompletionCosts = outputTokens * 0.002
		costs.TotalCosts = costs.PromptCosts + costs.CompletionCosts
	case openai.GPT3Dot5Turbo16K:
		costs.PromptCosts = inputTokens * 0.003
		costs.CompletionCosts = outputTokens * 0.004
		costs.TotalCosts = costs.PromptCosts + costs.CompletionCosts
	case openai.GPT4:
		costs.PromptCosts = inputTokens * 0.03
		costs.CompletionCosts = outputTokens * 0.06
		costs.TotalCosts = costs.PromptCosts + costs.CompletionCosts
	}

	return costs
}

func (chat Chat) RAG(ctx context.Context, prompt string) (*ChatCompletion, error) {

	embedding, err := chat.createEmbedding(ctx, prompt)
	if err != nil {
		return nil, err
	}

	sources, err := chat.db.SearchEmbedding(ctx, embedding)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(sources, func(i, j int) bool {
		return sources[i].Page < sources[j].Page
	})
	sort.SliceStable(sources, func(i, j int) bool {
		return sources[i].Source.String() < sources[j].Source.String()
	})

	var messages []openai.ChatCompletionMessage
	for _, result := range sources {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: result.Text,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	resp, err := chat.gpt.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       chat.Model,
			Temperature: float32(0),
			Messages:    messages,
			N:           1,
		},
	)

	if err != nil {
		return nil, err
	}

	return &ChatCompletion{
		Completion: resp.Choices[0].Message.Content,
		Sources:    sources,
		Costs:      chat.calculateCosts(resp.Usage),
	}, nil
}
