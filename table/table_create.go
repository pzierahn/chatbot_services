package table

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
)

// CreateTable creates a new table.
func (service *Service) CreateTable(ctx context.Context, req *pb.NewTable) (*pb.TableID, error) {
	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	var id string
	err = service.db.QueryRow(ctx,
		`INSERT INTO user_tables (user_id, name)
			VALUES ($1, $2)
            RETURNING id`,
		userId, req.Name).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &pb.TableID{Id: id}, nil
}
