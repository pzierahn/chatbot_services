package chat

import (
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
)

func (service *Service) getModel(name string) (llm.Completion, error) {
	switch name {
	case "gpt-4-1106-preview":
		return service.openai, nil
	case "gpt-4-turbo-preview":
		return service.openai, nil
	case "gpt-3.5-turbo-16k":
		return service.openai, nil
	case "gemini-pro":
		return service.vertex, nil
	}

	return nil, fmt.Errorf("unknown model: %v", name)
}
