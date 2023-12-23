package chat

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
)

func (service *Service) DeleteChatMessage(ctx context.Context, req *pb.MessageID) (*pb.MessageID, error) {
	userId, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	_, err = service.db.Exec(ctx, `
		DELETE FROM chat_messages
		WHERE user_id = $1
		  AND id = $2
	`, userId, req.Id)
	if err != nil {
		return nil, err
	}

	return req, nil
}
