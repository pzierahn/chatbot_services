package documents

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/pdf"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
	"strings"
	"sync"
	"sync/atomic"
)

type embeddingsBatch struct {
	userID uuid.UUID
	pages  []string
	stream pb.DocumentService_IndexServer
}

type document struct {
	id           string
	userId       string
	collectionId string
	filename     string
	path         string
	embeddings   []*embedding
}

type embedding struct {
	Page      uint32
	Text      string
	Embedding []float32
}

func (service *Service) getDocPages(ctx context.Context, path string) ([]string, error) {
	raw, err := service.storage.DownloadFile(bucket, path)
	if err != nil {
		return nil, err
	}

	return pdf.GetPagesFromBytes(ctx, raw)
}

func (service *Service) processEmbeddings(ctx context.Context, batch *embeddingsBatch) ([]*embedding, uint32, error) {
	totalPages := len(batch.pages)

	_ = batch.stream.Send(&pb.IndexProgress{
		TotalPages:     uint32(totalPages),
		ProcessedPages: 0,
	})

	var inputTokens uint32

	var wg sync.WaitGroup
	wg.Add(totalPages)

	var mu sync.Mutex
	var embeddings []*embedding
	var errs []error
	var processed uint32

	for inx, page := range batch.pages {
		go func(inx int, page string) {
			defer wg.Done()
			defer atomic.AddUint32(&processed, 1)

			page = strings.TrimSpace(page)
			if len(page) == 0 {
				return
			}

			resp, err := service.gpt.CreateEmbeddings(
				ctx,
				openai.EmbeddingRequestStrings{
					Model: embeddingsModel,
					Input: []string{page},
					User:  batch.userID.String(),
				},
			)

			mu.Lock()
			defer mu.Unlock()

			inputTokens += uint32(resp.Usage.PromptTokens)

			if err != nil {
				errs = append(errs, err)
				return
			}

			embeddings = append(embeddings, &embedding{
				Page:      uint32(inx),
				Text:      page,
				Embedding: resp.Data[0].Embedding,
			})

			_ = batch.stream.Send(&pb.IndexProgress{
				TotalPages:     uint32(totalPages),
				ProcessedPages: processed + 1,
			})
		}(inx, page)
	}

	wg.Wait()

	return embeddings, inputTokens, errors.Join(errs...)
}

func (service *Service) insertEmbeddings(ctx context.Context, doc *document) error {
	tx, err := service.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = service.db.Exec(
		ctx,
		`insert into documents (id, user_id, filename, path, collection_id)
			values ($1, $2, $3, $4, $5)`,
		doc.id,
		doc.userId,
		doc.filename,
		doc.path,
		doc.collectionId)
	if err != nil {
		return err
	}

	for _, fragment := range doc.embeddings {
		_, err = tx.Exec(ctx,
			`insert into document_embeddings (document_id, page, text, embedding)
				values ($1, $2, $3, $4)`,
			doc.id,
			fragment.Page,
			fragment.Text,
			pgvector.NewVector(fragment.Embedding))
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (service *Service) Index(doc *pb.Document, stream pb.DocumentService_IndexServer) error {

	ctx := stream.Context()
	userId, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return err
	}

	pages, err := service.getDocPages(ctx, doc.Path)
	if err != nil {
		return err
	}

	embeddings, inputTokens, err := service.processEmbeddings(ctx, &embeddingsBatch{
		userID: userId,
		pages:  pages,
		stream: stream,
	})
	if err != nil {
		service.storage.RemoveFile(bucket, []string{doc.Path})
		return err
	}

	err = service.insertEmbeddings(ctx, &document{
		id:           doc.Id,
		userId:       userId.String(),
		collectionId: doc.CollectionId,
		filename:     doc.Filename,
		path:         doc.Path,
		embeddings:   embeddings,
	})
	if err != nil {
		service.storage.RemoveFile(bucket, []string{doc.Path})
		return err
	}

	_, _ = service.account.CreateUsage(ctx, account.Usage{
		UserId: userId,
		Model:  embeddingsModel.String(),
		Input:  inputTokens,
	})

	return err
}
