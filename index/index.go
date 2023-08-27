package index

import (
	"braingain/database"
	"braingain/pdf"
	"context"
	"github.com/google/uuid"
	pb "github.com/qdrant/go-client/qdrant"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Index struct {
	collection string
	conn       *database.Client
	ai         *openai.Client
}

func (index Index) File(ctx context.Context, filename string) {
	log.Printf("Filename: %s\n", filepath.Base(filename))

	pages, err := pdf.ReadPages(ctx, filename)
	if err != nil {
		log.Fatalf("could not read pages: %v", err)
	}

	log.Printf("--> Pages: %d\n", len(pages))

	for inx, page := range pages {
		log.Printf("--> %d/%d\n", inx+1, len(pages))

		page = strings.TrimSpace(page)
		if len(page) == 0 {
			continue
		}

		resp, err := index.ai.CreateEmbeddings(
			context.Background(),
			openai.EmbeddingRequestStrings{
				Model: openai.AdaEmbeddingV2,
				Input: []string{page},
			},
		)

		if err != nil {
			log.Fatalf("could not create embeddings: %v", err)
		}

		err = index.conn.Upsert(ctx, database.Payload{
			Uuid:       uuid.NewString(),
			Collection: index.collection,
			Data:       resp.Data[0].Embedding,
			Metadata: map[string]*pb.Value{
				"filename": {
					Kind: &pb.Value_StringValue{
						StringValue: filepath.Base(filename),
					},
				},
				"page": {
					Kind: &pb.Value_IntegerValue{
						IntegerValue: int64(inx),
					},
				},
				"content": {
					Kind: &pb.Value_StringValue{
						StringValue: page,
					},
				},
			},
		})

		if err != nil {
			log.Fatalf("could not upsert points: %v", err)
		}
	}
}

func (index Index) Files(ctx context.Context, baseDir string) {
	_ = index.conn.DeleteCollection(ctx, index.collection)
	err := index.conn.CreateCollection(ctx, index.collection, 1536, pb.Distance_Cosine)
	if err != nil {
		log.Fatalf("could not create collection: %v", err)
	}

	// Read PDF files in baseDir
	files, err := os.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".pdf") {
			continue
		}

		pdfFile := baseDir + "/" + file.Name()
		index.File(ctx, pdfFile)
	}
}

func NewIndex(conn *database.Client, collection string) *Index {
	token := os.Getenv("OPENAI_API_KEY")
	ai := openai.NewClient(token)

	return &Index{
		collection: collection,
		conn:       conn,
		ai:         ai,
	}
}
