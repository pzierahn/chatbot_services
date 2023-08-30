package index

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/database"
	pb "github.com/qdrant/go-client/qdrant"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/net/html"
	"net/http"
	"path/filepath"
	"strings"
)

func (index Index) Web(ctx context.Context, url string) (err error) {
	// Get HTML from URL
	scrape, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = scrape.Body.Close() }()

	parts := make([]string, 0)

	// Extract text from HTML
	domDocTest := html.NewTokenizer(scrape.Body)
	previousStartTokenTest := domDocTest.Token()
loop:
	for {
		tt := domDocTest.Next()

		switch {
		case tt == html.ErrorToken:
			break loop
		case tt == html.StartTagToken:
			previousStartTokenTest = domDocTest.Token()
		case tt == html.TextToken:
			if previousStartTokenTest.Data == "script" {
				continue
			}

			text := strings.TrimSpace(html.UnescapeString(string(domDocTest.Text())))
			parts = append(parts, text)
		}
	}

	text := strings.Join(parts, "\n")

	resp, err := index.ai.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequestStrings{
			Model: openai.AdaEmbeddingV2,
			Input: []string{text},
		},
	)
	if err != nil {
		return err
	}

	err = index.conn.Upsert(ctx, database.Payload{
		Uuid:       uuid.NewString(),
		Collection: index.collection,
		Data:       resp.Data[0].Embedding,
		Metadata: map[string]*pb.Value{
			"url": {
				Kind: &pb.Value_StringValue{
					StringValue: filepath.Base(url),
				},
			},
			"content": {
				Kind: &pb.Value_StringValue{
					StringValue: text,
				},
			},
		},
	})

	return err
}
