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
	"strings"
)

func (service *Service) Index(req *pb.IndexJob, stream pb.DocumentService_IndexServer) error {

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

	var chunks []*pb.Chunk
	switch req.Document.Data.(type) {
	case *pb.DocumentMetadata_Web:
		_ = stream.Send(&pb.IndexProgress{
			Status: "Scraping webpage",
		})

		meta := req.Document.GetWeb()

		chunks, err = service.getWebChunks(ctx, meta)
	case *pb.DocumentMetadata_File:
		_ = stream.Send(&pb.IndexProgress{
			Status: "Extracting PDF pages",
		})

		meta := req.Document.GetFile()

		chunks, err = service.getPDFChunks(ctx, meta)
	default:
		return fmt.Errorf("unsupported metadata type")
	}

	if err != nil {
		return err
	}

	data := &document{
		userId:   userId,
		document: req,
		chunks:   chunks,
	}

	_ = stream.Send(&pb.IndexProgress{
		Status:   "Generating embeddings",
		Progress: 1.0 / 4.0,
	})
	embeddings, err := service.generateEmbeddings(ctx, userId, data)
	if err != nil {
		return err
	}

	_ = stream.Send(&pb.IndexProgress{
		Status:   "Inserting into database",
		Progress: 2.0 / 4.0,
	})
	err = service.insertIntoDB(ctx, data)
	if err != nil {
		return err
	}

	_ = stream.Send(&pb.IndexProgress{
		Status:   "Inserting into vector database",
		Progress: 3.0 / 4.0,
	})
	err = service.insertEmbeddings(data, embeddings)
	if err != nil {
		return err
	}

	_ = stream.Send(&pb.IndexProgress{
		Status:   "Success",
		Progress: 1.0,
	})

	return nil
}

func (service *Service) getWebChunks(ctx context.Context, meta *pb.Webpage) ([]*pb.Chunk, error) {
	text, err := web.Scrape(ctx, meta.Url)
	if err != nil {
		return nil, err
	}

	var inx uint32
	var chunks []*pb.Chunk

	for chunk := 0; chunk < len(text)/6144; chunk++ {

		start := max(chunk*6144-200, 0)
		end := min((chunk+1)*6144+200, len(text))
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

	chunks := make([]*pb.Chunk, len(pages))
	for inx, page := range pages {
		chunks[inx] = &pb.Chunk{
			Id:    uuid.NewString(),
			Text:  page,
			Index: uint32(inx),
		}
	}

	return chunks, nil
}
