package chat_v5

import (
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/auth"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/llm/openai"
	pb "github.com/pzierahn/chatbot_services/proto"
)

type Service struct {
	pb.UnimplementedChatServiceServer
	models []llm.Chat
	auth   auth.Service
	db     *datastore.Service
}

// getModel returns the llm.Chat that provides the given model.
func (service *Service) getModel(name string) (llm.Chat, error) {
	for _, model := range service.models {
		if model.ProvidesModel(name) {
			return model, nil
		}
	}

	return nil, fmt.Errorf("model not found: %s", name)
}

// Verify checks the user's authentication and funding status.
func (service *Service) Verify(ctx context.Context) (string, error) {
	userId, err := service.auth.Verify(ctx)
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

func New() (*Service, error) {
	//uri := os.Getenv("CHATBOT_MONGODB_URI")
	//if uri == "" {
	//	log.Fatal("MONGODB_URI is not set")
	//}

	ctx := context.Background()
	db, err := datastore.New(ctx)
	if err != nil {
		return nil, err
	}

	fakeAuth, _ := auth.WithInsecure()

	openaiClient, err := openai.New()
	if err != nil {
		return nil, err
	}

	models := []llm.Chat{
		openaiClient,
	}

	return &Service{
		auth:   fakeAuth,
		db:     db,
		models: models,
	}, nil
}
