package documents

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/search"
	"time"
)

func (service *Service) addToSearchIndex(ctx context.Context, doc *datastore.Document) error {
	var vectors []*search.Fragment

	for _, fragment := range doc.Content {
		if fragment.Id == uuid.Nil {
			return fmt.Errorf("fragment id is empty")
		}

		if fragment.Text == "" {
			continue
		}

		vectors = append(vectors, &search.Fragment{
			Id:           fragment.Id.String(),
			Text:         fragment.Text,
			UserId:       doc.UserId,
			DocumentId:   doc.Id.String(),
			CollectionId: doc.CollectionId.String(),
			Position:     fragment.Position,
		})
	}

	usage, err := service.SearchIndex.Upsert(ctx, vectors)
	if err != nil {
		return err
	}

	_ = service.Database.InsertModelUsage(ctx, &datastore.ModelUsage{
		Id:          uuid.New(),
		UserId:      doc.UserId,
		Timestamp:   time.Now(),
		ModelId:     usage.ModelId,
		InputTokens: usage.Tokens,
	})

	return nil
}
