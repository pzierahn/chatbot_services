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

	data := document{
		userId:    userId,
		document:  req,
		chunkMeta: chunks,
	}

	_ = stream.Send(&pb.IndexProgress{
		Status:   "Inserting into database",
		Progress: 0.77,
	})

	err = service.insertIntoDB(ctx, data)
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

	return []*pb.Chunk{
		{
			Id:       uuid.NewString(),
			Text:     text,
			Metadata: &pb.Chunk_Web{Web: &pb.WebpageChunkMetadata{}},
		},
	}, nil
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
			Id:   uuid.NewString(),
			Text: page,
			Metadata: &pb.Chunk_Doc{
				Doc: &pb.FileChunkMetadata{
					Page: uint32(inx),
				},
			},
		}
	}

	return chunks, nil
}
