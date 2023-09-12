package index

import (
	"context"
	"github.com/pzierahn/braingain/database"
	"github.com/sashabaranov/go-openai"
	"strings"
)

func (index Index) GetPagesWithEmbeddings(ctx context.Context, pages []string) ([]*database.PageEmbedding, error) {
	var embeddings []*database.PageEmbedding

	for inx, page := range pages {
		page = strings.TrimSpace(page)
		if len(page) == 0 {
			continue
		}

		resp, err := index.GPT.CreateEmbeddings(
			ctx,
			openai.EmbeddingRequestStrings{
				Model: openai.AdaEmbeddingV2,
				Input: []string{page},
			},
		)
		if err != nil {
			return nil, err
		}

		embeddings = append(embeddings, &database.PageEmbedding{
			Page:      inx,
			Text:      page,
			Embedding: resp.Data[0].Embedding,
		})
	}

	return embeddings, nil
}
