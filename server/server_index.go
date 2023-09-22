package server

import (
	"errors"
	"github.com/google/uuid"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/database"
	"github.com/pzierahn/brainboost/pdf"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
	"log"
	"strings"
	"sync"
	"sync/atomic"
)

const bucket = "documents"

type Progress struct {
	TotalPages   int
	FinishedPage int
}

func (server *Server) IndexDocument(doc *pb.Document, stream pb.Brainboost_IndexDocumentServer) error {

	uid, err := auth.ValidateToken(stream.Context())
	if err != nil {
		return err
	}

	log.Printf("IndexDocument: %+v", doc)

	ctx := stream.Context()

	raw, err := server.storage.DownloadFile(bucket, doc.Path)
	if err != nil {
		return err
	}

	pages, err := pdf.GetPagesFromBytes(ctx, raw)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(pages))

	var mu sync.Mutex
	var embeddings []*database.PageEmbedding
	var errs []error
	var inputTokens int
	var processed uint32

	_ = stream.Send(&pb.IndexProgress{
		TotalPages:     uint32(len(pages)),
		ProcessedPages: 0,
	})

	for inx, page := range pages {
		go func(inx int, page string) {
			defer wg.Done()
			defer atomic.AddUint32(&processed, 1)

			page = strings.TrimSpace(page)
			if len(page) == 0 {
				return
			}

			resp, err := server.gpt.CreateEmbeddings(
				ctx,
				openai.EmbeddingRequestStrings{
					Model: embeddingsModel,
					Input: []string{page},
					User:  uid.String(),
				},
			)

			mu.Lock()
			defer mu.Unlock()

			inputTokens += resp.Usage.PromptTokens

			if err != nil {
				errs = append(errs, err)
				return
			}

			embeddings = append(embeddings, &database.PageEmbedding{
				Page:      inx,
				Text:      page,
				Embedding: resp.Data[0].Embedding,
			})

			_ = stream.Send(&pb.IndexProgress{
				TotalPages:     uint32(len(pages)),
				ProcessedPages: processed + 1,
			})
		}(inx, page)
	}

	wg.Wait()

	err = errors.Join(errs...)
	if err != nil {
		server.storage.RemoveFile(bucket, []string{doc.Path})
		return err
	}

	_, err = server.db.UpsertDocument(ctx, database.Document{
		UserId:     uid.String(),
		Collection: uuid.MustParse(doc.CollectionId),
		Filename:   doc.Filename,
		Path:       doc.Path,
		Pages:      embeddings,
	})
	if err != nil {
		server.storage.RemoveFile(bucket, []string{doc.Path})
		return err
	}

	log.Printf("Indexing done: %v", doc.Filename)

	_, _ = server.db.CreateUsage(ctx, database.Usage{
		UID:   uid.String(),
		Model: embeddingsModel.String(),
		Input: inputTokens,
	})

	return err
}
