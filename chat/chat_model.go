package chat

import (
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
)

func (service *Service) getModel(name string) (llm.Completion, error) {
	for _, model := range service.models {
		if model.ProvideModel(name) {
			return model, nil
		}
	}

	return nil, fmt.Errorf("unknown model: %v", name)
}
