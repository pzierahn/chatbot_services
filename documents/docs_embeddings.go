package documents

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/search"
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

	return service.SearchIndex.Upsert(ctx, vectors)
}
