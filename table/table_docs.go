package table

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/llm/bedrock"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func docToText(doc *pb.Document) string {
	var text string

	for _, chunk := range doc.Chunks {
		text += chunk.Text
	}

	return text
}

func (service *Service) storeRows(ctx context.Context, userId, tableId string, rows []*pb.Row) error {
	for _, row := range rows {
		var rowId string
		err := service.db.QueryRow(ctx,
			`INSERT INTO user_table_rows (user_id, table_id, document_id)
				VALUES ($1, $2, $3)
				RETURNING id`,
			userId, tableId, row.DocumentId).Scan(&rowId)
		if err != nil {
			return err
		}

		for _, cell := range row.Cells {
			_, err = service.db.Exec(ctx,
				`INSERT INTO user_table_cells (user_id, table_id, row_id, column_id, value)
					VALUES ($1, $2, $3, $4, $5)`,
				userId, tableId, rowId, cell.ColumnId, cell.Value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (service *Service) AddDocumentsToTable(ctx context.Context, req *pb.DocumentsToTable) (*emptypb.Empty, error) {
	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	columns, err := service.getColumns(ctx, userId, req.TableId)
	if err != nil {
		return nil, err
	}

	var rows []*pb.Row

	for _, docID := range req.DocumentIds {
		doc, err := service.document.Get(ctx, &pb.DocumentID{Id: docID})
		if err != nil {
			return nil, err
		}

		text := docToText(doc)

		row := &pb.Row{
			DocumentId: docID,
		}

		for _, col := range columns {
			log.Printf("Generating completion for column %s", col.Id)
			resp, err := service.agent.GenerateCompletion(ctx, &llm.GenerateRequest{
				SystemPrompt: "Return only the requested information, without any additional words or explanation. " +
					"Don't use introductory words. " +
					"The answer should be as short as possible. " +
					"Don't repeat parts of the question or prompt.",
				Messages: []*llm.Message{{
					Type: llm.MessageTypeUser,
					Text: text + "\n\n" + col.GenerationPrompt,
				}},
				Model:       bedrock.ClaudeSonnet,
				MaxTokens:   100,
				TopP:        0,
				Temperature: 0,
				UserId:      userId,
			})
			if err != nil {
				return nil, err
			}

			row.Cells = append(row.Cells, &pb.Cell{
				ColumnId: col.Id,
				Value:    resp.Text,
			})
		}

		rows = append(rows, row)
	}

	err = service.storeRows(ctx, userId, req.TableId, rows)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
