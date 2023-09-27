package chat

import (
	"context"
	pb "github.com/pzierahn/brainboost/proto"
	"sort"
	"strings"
)

func (service *Service) getSourceFromDB(ctx context.Context, prompt *pb.Prompt) (*chatContext, error) {

	query := &pb.SearchQuery{
		CollectionID: prompt.CollectionID,
		Query:        prompt.Prompt,
		Limit:        prompt.Options.Limit,
		Threshold:    prompt.Options.Threshold,
	}

	results, err := service.docs.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	sort.Slice(results.Items, func(i, j int) bool {
		return results.Items[i].Page > results.Items[j].Page
	})

	bg := chatContext{}

	filename := make(map[string]string)
	pages := make(map[string][]uint32)
	scores := make(map[string][]float32)
	text := make(map[string][]string)

	for _, doc := range results.Items {
		filename[doc.Id] = doc.Filename
		pages[doc.Id] = append(pages[doc.Id], doc.Page)
		scores[doc.Id] = append(scores[doc.Id], doc.Score)
		text[doc.Id] = append(text[doc.Id], doc.Content)
	}

	for id := range filename {
		bg.fragments = append(bg.fragments, strings.Join(text[id], "\n"))
		bg.docs = append(bg.docs, &pb.ChatMessage_Document{
			Id:       id,
			Filename: filename[id],
			Pages:    pages[id],
			Scores:   scores[id],
		})
	}

	return &bg, nil
}
