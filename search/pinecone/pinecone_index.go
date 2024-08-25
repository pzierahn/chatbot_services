package pinecone_search

import (
	"context"
	"fmt"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/pzierahn/chatbot_services/search"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
)

func (db *Search) getIndexConnection(ctx context.Context) (*pinecone.IndexConnection, error) {
	idx, err := db.conn.DescribeIndex(ctx, db.namespace)
	if err != nil {
		return nil, err
	}

	idxConnection, err := db.conn.Index(pinecone.NewIndexConnParams{Host: idx.Host})
	if err != nil {
		return nil, err
	}

	return idxConnection, nil
}

func (db *Search) Upsert(ctx context.Context, fragments []*search.Fragment) (*search.Usage, error) {

	embedded, err := db.fastEmbedding.CreateEmbeddings(ctx, fragments)
	if err != nil {
		return nil, err
	}

	var vectors []*pinecone.Vector

	for idx := range fragments {
		fragment := fragments[idx]

		metadata, err := structpb.NewStruct(map[string]any{
			search.PayloadFragmentId:   fragment.Id,
			search.PayloadDocumentId:   fragment.DocumentId,
			search.PayloadUserId:       fragment.UserId,
			search.PayloadCollectionId: fragment.CollectionId,
			search.PayloadText:         fragment.Text,
			search.PayloadPosition:     fragment.Position,
		})
		if err != nil {
			return nil, err
		}

		vectorId := fmt.Sprintf("%s#%s#%s", fragment.CollectionId, fragment.DocumentId, fragment.Id)

		vectors = append(vectors, &pinecone.Vector{
			Id:       vectorId,
			Values:   embedded.Embeddings[fragment.Id],
			Metadata: metadata,
		})
	}

	idxConnection, err := db.getIndexConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = idxConnection.Close() }()

	count, err := idxConnection.UpsertVectors(ctx, vectors)
	if err != nil {
		return nil, err
	} else {
		log.Printf("Successfully upserted %d vector(s)!\n", count)
	}

	return &embedded.Usage, nil
}
