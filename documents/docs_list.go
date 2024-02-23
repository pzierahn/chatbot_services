package documents

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/utils"
	"strings"
)

func (service *Service) List(ctx context.Context, req *pb.DocumentFilter) (*pb.DocumentList, error) {

	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := service.db.Query(ctx,
		`SELECT id, metadata
		FROM documents
		WHERE
		    user_id = $1 AND
		    collection_id = $2::uuid`,
		userId, req.CollectionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documents := &pb.DocumentList{
		Items: make(map[string]*pb.DocumentMetadata),
	}

	for rows.Next() {
		var (
			docId string
			meta  DocumentMeta
		)

		err = rows.Scan(
			&docId,
			&meta,
		)
		if err != nil {
			return nil, err
		}

		proto := meta.toProto()
		title := strings.ToLower(utils.GetDocumentTitle(proto))
		query := strings.ToLower(req.Query)

		if strings.Contains(title, query) {
			documents.Items[docId] = proto
		}
	}

	return documents, nil
}
