package documents

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/vectordb"
)

func (service *Service) addToSearchIndex(ctx context.Context, doc *datastore.Document) error {
	var vectors []*vectordb.Fragment

	for _, fragment := range doc.Content {
		if fragment.Id == uuid.Nil {
			return fmt.Errorf("fragment id is empty")
		}

		vectors = append(vectors, &vectordb.Fragment{
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
