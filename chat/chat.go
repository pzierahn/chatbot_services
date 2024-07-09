package chat

import (
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/vectordb"
)

type Service struct {
	pb.UnimplementedChatServiceServer
	Models   []llm.Chat
	Auth     account.Verifier
	Database *datastore.Service
	Search   vectordb.DB
}

// getModel returns the llm.Chat that provides the given model.
func (service *Service) getModel(name string) (llm.Chat, error) {
	for _, model := range service.Models {
		if model.ProvidesModel(name) {
			return model, nil
		}
	}

	return nil, fmt.Errorf("model not found: %s", name)
}

// Verify checks the user's authentication and funding status.
func (service *Service) Verify(ctx context.Context) (string, error) {
	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return "", err
	}

	funding, err := service.hasFunding(ctx, userId)
	if err != nil {
		return "", err
	}

	if !funding {
		return "", account.NoFundingError()
	}

	return userId, nil
}

func (service *Service) hasFunding(ctx context.Context, userId string) (bool, error) {
	// TODO: Check if the user has enough funds to chat.
	return true, nil
}
