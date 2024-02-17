package documents

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/pdf"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/web"
	"github.com/sashabaranov/go-openai"
	"io"
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

	var chunks []*pb.Chunk
	switch req.Document.Data.(type) {
	case *pb.DocumentMetadata_Web:
		_ = stream.Send(&pb.IndexProgress{
			Status: "Scraping webpage",
		})

		meta := req.Document.GetWeb()
		chunks, err = service.getWebChunks(ctx, userId, meta)
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

func (service *Service) getWebChunks(ctx context.Context, userId string, meta *pb.Webpage) ([]*pb.Chunk, error) {
	text, err := web.Scrape(ctx, meta.Url)
	if err != nil {
		return nil, err
	}

	resp, err := service.LLM.GenerateCompletion(ctx, &llm.GenerateRequest{
		Messages: []*llm.Message{
			{
				Type: llm.MessageTypeUser,
				Text: text,
			},
			{
				Type: llm.MessageTypeUser,
				Text: "Split this text into chunks, operate each chunk by %%%%%%%%%%",
			},
		},
		Model:  openai.GPT4TurboPreview,
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	var chunks []*pb.Chunk
	for inx, chunk := range strings.Split(resp.Text, "%%%%%%%%%%") {
		chunks = append(chunks, &pb.Chunk{
			Id:    uuid.NewString(),
			Text:  strings.TrimSpace(chunk),
			Index: uint32(inx),
		})

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
