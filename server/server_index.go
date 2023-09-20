package server

import (
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/auth"
	"github.com/pzierahn/braingain/index"
	pb "github.com/pzierahn/braingain/proto"
	"log"
)

func (server *Server) IndexDocument(ref *pb.StorageRef, stream pb.Braingain_IndexDocumentServer) error {

	uid, err := auth.ValidateToken(stream.Context())
	if err != nil {
		return err
	}

	log.Printf("Indexing: %v", ref)

	docId := index.DocumentId{
		UserId:     uid.String(),
		Collection: uuid.MustParse(ref.CollectionId),
		DocId:      uuid.MustParse(ref.Id),
		Filename:   ref.Filename,
	}

	byt, err := server.index.Download(docId)
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

	_, err = server.index.Process(ctx, docId, byt, progress)
	log.Printf("Indexing done: %v", err)

	return err
}
