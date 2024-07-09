package chat

import (
	"fmt"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/search"
)

type Service struct {
	pb.UnimplementedChatServiceServer
	Models   []llm.Chat
	Auth     account.Verifier
	Database *datastore.Service
	Search   search.DB
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
