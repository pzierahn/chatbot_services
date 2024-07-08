package documents

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/pdf"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/web"
	"io"
	"strings"
)

func (service *Service) Index(req *pb.IndexJob, stream pb.DocumentService_IndexServer) error {
	ctx := stream.Context()

	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return err
	}

	var documentId uuid.UUID
	if req.Id == "" {
		documentId = uuid.New()
	} else {
		documentId, err = uuid.Parse(req.Id)
		if err != nil {
			return err
		}
	}

	collectionId, err := uuid.Parse(req.CollectionId)
	if err != nil {
		return err
	}

	data := &datastore.Document{
		Id:           documentId,
		UserId:       userId,
		CollectionId: collectionId,
		Name:         "",
		Type:         "",
		Source:       "",
	}

	switch req.Document.Data.(type) {
	case *pb.DocumentMetadata_Web:
		_ = stream.Send(&pb.IndexProgress{
			Status: "Scraping webpage",
		})

		meta := req.Document.GetWeb()
		data.Type = datastore.DocumentTypeWeb
		data.Name = meta.Title
		data.Source = meta.Url
		data.Content, err = service.getWebChunks(ctx, meta)
	case *pb.DocumentMetadata_File:
		_ = stream.Send(&pb.IndexProgress{
			Status: "Extracting PDF pages",
		})
		meta := req.Document.GetFile()
		data.Type = datastore.DocumentTypePDF
		data.Name = meta.Filename
		data.Source = meta.Path
		data.Content, err = service.getPDFChunks(ctx, meta)
	default:
		return fmt.Errorf("unsupported metadata type")
	}

	if err != nil {
		return err
	}

	_ = stream.Send(&pb.IndexProgress{
		Status:   "Inserting into search database",
		Progress: 1.0 / 3.0,
	})
	err = service.addToSearchIndex(ctx, data)
	if err != nil {
		return err
	}

	_ = stream.Send(&pb.IndexProgress{
		Status:   "Inserting into database",
		Progress: 2.0 / 3.0,
	})
	err = service.Database.StoreDocument(ctx, data)
	if err != nil {
		return err
	}

	_ = stream.Send(&pb.IndexProgress{
		Status:   "Success",
		Progress: 1.0,
	})

	return nil
}

func (service *Service) getWebChunks(ctx context.Context, meta *pb.Webpage) ([]*datastore.DocumentChunk, error) {
	text, err := web.Scrape(ctx, meta.Url)
	if err != nil {
		return nil, err
	}

	var inx uint32
	var chunks []*datastore.DocumentChunk

	for chunk := 0; chunk < len(text)/6144; chunk++ {

		start := max(chunk*6144-200, 0)
		end := min((chunk+1)*6144+200, len(text))
		fragment := text[start:end]

		chunks = append(chunks, &datastore.DocumentChunk{
			Id:       uuid.New(),
			Text:     strings.TrimSpace(fragment),
			Position: inx,
		})

		inx++
	}

	return chunks, nil
}

func (service *Service) getPDFChunks(ctx context.Context, meta *pb.File) ([]*datastore.DocumentChunk, error) {

	obj := service.Storage.Object(meta.Path)
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

	chunks := make([]*datastore.DocumentChunk, len(pages))
	for inx, page := range pages {
		chunks[inx] = &datastore.DocumentChunk{
			Id:       uuid.New(),
			Text:     strings.TrimSpace(page),
			Position: uint32(inx),
		}
	}

	return chunks, nil
}
