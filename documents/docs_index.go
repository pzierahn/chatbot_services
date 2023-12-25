package documents

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/pdf"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/vectordb"
	"io"
	"strings"
	"sync"
	"sync/atomic"
)

type embeddingsBatch struct {
	userId string
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

	obj := service.storage.Object(path)
	read, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = read.Close() }()

	raw, err := io.ReadAll(read)
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

			resp, err := service.embeddings.CreateEmbeddings(ctx, &llm.EmbeddingRequest{
				Input:  page,
				UserId: batch.userId,
			})

			mu.Lock()
			defer mu.Unlock()

			inputTokens += uint32(resp.Tokens)

			if err != nil {
				errs = append(errs, err)
				return
			}

			embeddings = append(embeddings, &embedding{
				Page:      uint32(inx),
				Text:      page,
				Embedding: resp.Data,
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

	var vectors []*vectordb.Vector

	for _, fragment := range doc.embeddings {
		chunkId := uuid.NewString()

		_, err = tx.Exec(ctx,
			`insert into document_chunks (id, document_id, page, text)
				values ($1, $2, $3, $4)`,
			chunkId,
			doc.id,
			fragment.Page,
			fragment.Text,
		)
		if err != nil {
			return err
		}

		vectors = append(vectors, &vectordb.Vector{
			Id:           chunkId,
			DocumentId:   doc.id,
			CollectionId: doc.collectionId,
			UserId:       doc.userId,
			Filename:     doc.filename,
			Text:         fragment.Text,
			Page:         fragment.Page,
			Vector:       fragment.Embedding,
		})
	}

	err = service.vectorDB.Upsert(vectors)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (service *Service) Index(doc *pb.Document, stream pb.DocumentService_IndexServer) error {

	ctx := stream.Context()
	userId, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return err
	}

	funding, err := service.account.HasFunding(ctx)
	if err != nil {
		return err
	}

	if !funding {
		return account.NoFundingError()
	}

	pages, err := service.getDocPages(ctx, doc.Path)
	if err != nil {
		return err
	}

	obj := service.storage.Object(doc.Path)

	embeddings, inputTokens, err := service.processEmbeddings(ctx, &embeddingsBatch{
		userId: userId,
		pages:  pages,
		stream: stream,
	})
	if err != nil {
		_ = obj.Delete(ctx)
		return err
	}

	err = service.insertEmbeddings(ctx, &document{
		id:           doc.Id,
		userId:       userId,
		collectionId: doc.CollectionId,
		filename:     doc.Filename,
		path:         doc.Path,
		embeddings:   embeddings,
	})
	if err != nil {
		_ = obj.Delete(ctx)
		return err
	}

	_, _ = service.account.CreateUsage(ctx, account.Usage{
		UserId: userId,
		Model:  embeddingsModel.String(),
		Input:  inputTokens,
	})

	return err
}
