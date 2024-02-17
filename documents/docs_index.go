package documents

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/pdf"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/web"
	"io"
	"net/url"
	"strings"
)

func (service *Service) IndexDocument(req *pb.IndexJob, stream pb.DocumentService_IndexDocumentServer) error {

	ctx := stream.Context()

	userId, err := service.auth.Verify(ctx)
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

	if req.Id == "" {
		req.Id = uuid.NewString()
	}

	var title string
	var chunks []*pb.Chunk
	switch req.Document.Data.(type) {
	case *pb.DocumentMetadata_Web:
		_ = stream.Send(&pb.IndexProgress{
			Status: "Scraping webpage",
		})

		meta := req.Document.GetWeb()
		metaUrl, err := url.Parse(meta.Url)
		if err != nil {
			return err
		}

		title = metaUrl.Host + metaUrl.Path

		chunks, err = service.getWebChunks(ctx, meta)
	case *pb.DocumentMetadata_File:
		_ = stream.Send(&pb.IndexProgress{
			Status: "Extracting PDF pages",
		})

		meta := req.Document.GetFile()
		title = meta.Filename

		chunks, err = service.getPDFChunks(ctx, meta)
	default:
		return fmt.Errorf("unsupported metadata type")
	}

	if err != nil {
		return err
	}

	data := &document{
		userId:   userId,
		title:    title,
		document: req,
		chunks:   chunks,
	}

	_ = stream.Send(&pb.IndexProgress{
		Status: "Generating embeddings",
	})
	embeddings, err := service.generateEmbeddings(ctx, userId, data)
	if err != nil {
		return err
	}

	_ = stream.Send(&pb.IndexProgress{
		Status: "Inserting into database",
	})
	err = service.insertIntoDB(ctx, data)
	if err != nil {
		return err
	}

	_ = stream.Send(&pb.IndexProgress{
		Status: "Inserting into vector database",
	})
	err = service.insertEmbeddings(data, embeddings)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) getWebChunks(ctx context.Context, meta *pb.Webpage) ([]*pb.Chunk, error) {
	text, err := web.Scrape(ctx, meta.Url)
	if err != nil {
		return nil, err
	}

	var inx uint32
	var chunks []*pb.Chunk

	for chunk := 0; chunk < len(text)/3072; chunk++ {

		start := max(chunk*3072-100, 0)
		end := min((chunk+1)*3072+100, len(text))
		fragment := text[start:end]

		chunks = append(chunks, &pb.Chunk{
			Id:    uuid.NewString(),
			Text:  strings.TrimSpace(fragment),
			Index: inx,
		})

		inx++
	}

	return chunks, nil
}

func (service *Service) getPDFChunks(ctx context.Context, meta *pb.File) ([]*pb.Chunk, error) {

	obj := service.storage.Object(meta.Path)
	read, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = read.Close() }()

	raw, err := io.ReadAll(read)
	if err != nil {
		return nil, err
	}

	pages, err := pdf.GetPagesFromBytes(ctx, raw)
	if err != nil {
		return nil, err
	}

	chunks := make([]*pb.Chunk, 0, len(pages))
	for inx, page := range pages {
		chunks[inx] = &pb.Chunk{
			Id:    uuid.NewString(),
			Text:  page,
			Index: uint32(inx),
		}
	}

	return chunks, nil
}
