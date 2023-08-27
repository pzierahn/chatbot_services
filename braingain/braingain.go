package braingain

import (
	"context"
	"github.com/pzierahn/braingain/database"
	"github.com/sashabaranov/go-openai"
	"sort"
)

const (
	collection = "DeSys"
)

type Source struct {
	Id       string
	Score    float32
	Filename string
	Page     int
	Content  string
}

type ChatCompletion struct {
	Completion string
	Sources    []Source
}

type Chat struct {
	db  *database.Client
	gpt *openai.Client
}

func NewChat(db *database.Client, gpt *openai.Client) *Chat {
	return &Chat{
		db:  db,
		gpt: gpt,
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

func (chat Chat) RAG(ctx context.Context, prompt string) (*ChatCompletion, error) {

	embedding, err := chat.createEmbedding(ctx, prompt)
	if err != nil {
		return nil, err
	}

	searchResponse, err := chat.db.SearchEmbedding(ctx, collection, embedding)
	if err != nil {
		return nil, err
	}

	sources := make([]Source, 0)

	for _, hit := range searchResponse.Result {
		filename := hit.Payload["filename"].GetStringValue()
		page := int(hit.Payload["page"].GetIntegerValue()) + 1

		sources = append(sources, Source{
			Id:       hit.Id.GetUuid(),
			Score:    hit.Score,
			Filename: filename,
			Page:     page,
			Content:  hit.Payload["content"].GetStringValue(),
		})
	}

	sort.SliceStable(sources, func(i, j int) bool {
		return sources[i].Page < sources[j].Page
	})
	sort.SliceStable(sources, func(i, j int) bool {
		return sources[i].Filename < sources[j].Filename
	})

	var messages []openai.ChatCompletionMessage
	for _, result := range sources {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: result.Content,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	resp, err := chat.gpt.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			//Model: openai.GPT3Dot5Turbo16K,
			Model:       openai.GPT4,
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
	}, nil
}
