package server

import (
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/index"
	pb "github.com/pzierahn/braingain/proto"
	"log"
)

func (server *Server) IndexDocument(ref *pb.StorageRef, stream pb.Braingain_IndexDocumentServer) error {

	log.Printf("Indexing: %v", ref)

	source := index.Index{
		DB:      server.db,
		GPT:     server.gpt,
		Storage: server.storage,
	}

	docId := index.DocumentId{
		UserId:     patrick,
		Collection: uuid.MustParse(ref.Collection),
		DocId:      uuid.MustParse(ref.Id),
		Filename:   ref.Filename,
	}

	byt, err := source.Download(docId)
	if err != nil {
		return err
	}

	ctx := stream.Context()

	progress := make(chan index.Progress)
	defer close(progress)

	go func() {
		var processed uint32
		for p := range progress {
			processed += 1
			log.Printf("Progress: %v/%v", processed, p.TotalPages)

			_ = stream.Send(&pb.IndexProgress{
				TotalPages:     uint32(p.TotalPages),
				ProcessedPages: processed,
			})
		}
	}()

	_, err = source.Process(ctx, docId, byt, progress)
	log.Printf("Indexing done: %v", err)

	return err
}
