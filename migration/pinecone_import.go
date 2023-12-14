package migration

import (
	"context"
	"crypto/tls"
	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"os"
)

func PineconeImport(ctx context.Context) {
	config := &tls.Config{}

	userMapping := GetUserIdMapping()

	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", os.Getenv("PINECONE_KEY"))
	target := os.Getenv("PINECONE_URL")

	log.Printf("connecting to %v", target)

	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithAuthority(target),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	client := pinecone_grpc.NewVectorServiceClient(conn)

	docs := ExportDocumentsMeta(ctx)
	log.Printf("docs: %v", len(docs))

	for inx, doc := range docs {
		log.Printf("doc: %v (%d/%d)", doc.Id, inx, len(docs))

		embeddings := ExportDocumentVectors(ctx, doc)
		log.Printf("embeddings: %v", len(embeddings))

		if len(embeddings) == 0 {
			continue
		}

		var vectors []*pinecone_grpc.Vector

		newUserId, ok := userMapping[doc.UserId]
		if !ok {
			log.Fatalf("user not found: %v", doc.UserId)
		}

		for _, embedding := range embeddings {
			meta := &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"documentId":   {Kind: &structpb.Value_StringValue{StringValue: embedding.DocumentId}},
					"collectionId": {Kind: &structpb.Value_StringValue{StringValue: doc.CollectionId}},
					"userId":       {Kind: &structpb.Value_StringValue{StringValue: newUserId}},
					"filename":     {Kind: &structpb.Value_StringValue{StringValue: doc.Filename}},
					"text":         {Kind: &structpb.Value_StringValue{StringValue: embedding.Text}},
					"page":         {Kind: &structpb.Value_NumberValue{NumberValue: float64(embedding.Page)}},
				},
			}

			vectors = append(vectors, &pinecone_grpc.Vector{
				Id:       embedding.Id,
				Values:   embedding.Embedding,
				Metadata: meta,
			})
		}

		for iny := 0; iny < len(vectors); iny += 50 {
			end := min(iny+50, len(vectors))

			_, err = client.Upsert(ctx, &pinecone_grpc.UpsertRequest{
				Vectors:   vectors[iny:end],
				Namespace: "documents",
			})

			if err != nil {
				log.Fatalf("upsert error: %v", err)
			}
		}
	}
}
