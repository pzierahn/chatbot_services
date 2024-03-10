package table

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"sort"
)

func (service *Service) getColumns(ctx context.Context, userId, tableId string) ([]*pb.Column, error) {
	rows, err := service.db.Query(ctx,
		`SELECT id, name, generation_prompt 
			FROM user_table_columns
			WHERE user_id = $1 AND
			      table_id = $2`, userId, tableId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns := make([]*pb.Column, 0)
	for rows.Next() {
		col := new(pb.Column)
		err = rows.Scan(&col.Id, &col.Name, &col.GenerationPrompt)
		if err != nil {
			return nil, err
		}

		columns = append(columns, col)
	}

	return columns, nil
}

func (service *Service) getDocumentIds(ctx context.Context, userId, tableId string) (map[string]string, error) {
	rows, err := service.db.Query(ctx,
		`SELECT id, document_id 
			FROM user_table_rows
			WHERE user_id = $1 AND
			      table_id = $2`, userId, tableId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documentIds := make(map[string]string)
	for rows.Next() {
		var id, documentId string
		err = rows.Scan(&id, &documentId)
		if err != nil {
			return nil, err
		}

		documentIds[id] = documentId
	}

	return documentIds, nil
}

func (service *Service) getRows(ctx context.Context, userId, tableId string) ([]*pb.Row, error) {
	documentIds, err := service.getDocumentIds(ctx, userId, tableId)
	if err != nil {
		return nil, err
	}

	rows, err := service.db.Query(ctx,
		`SELECT column_id, row_id, value 
			FROM user_table_cells
			WHERE user_id = $1 AND
			      table_id = $2`, userId, tableId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data := make(map[string][]*pb.Cell)

	for rows.Next() {
		var columnId, rowId, value string
		err = rows.Scan(&columnId, &rowId, &value)
		if err != nil {
			return nil, err
		}

		docId, ok := documentIds[rowId]
		if !ok {
			continue
		}

		data[docId] = append(data[rowId], &pb.Cell{
			ColumnId: columnId,
			Value:    value,
		})
	}

	var dataRows []*pb.Row
	for docId, cells := range data {
		dataRows = append(dataRows, &pb.Row{
			DocumentId: docId,
			Cells:      cells,
		})
	}

	sort.Slice(dataRows, func(i, j int) bool {
		return dataRows[i].DocumentId < dataRows[j].DocumentId
	})

	return dataRows, nil
}

func (service *Service) GetTable(ctx context.Context, req *pb.TableID) (*pb.Table, error) {
	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	columns, err := service.getColumns(ctx, userId, req.Id)
	if err != nil {
		return nil, err
	}

	rows, err := service.getRows(ctx, userId, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Table{
		Id:      req.Id,
		Columns: columns,
		Rows:    rows,
	}, nil
}
