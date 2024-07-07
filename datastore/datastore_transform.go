package datastore

import (
	"encoding/json"
	"github.com/pzierahn/chatbot_services/llm"
)

// ToLLMMessage converts a slice of datastore messages to a slice of LLM messages.
func ToLLMMessage(message *Message) (*llm.Message, error) {
	result := &llm.Message{
		Role:          message.Role,
		Content:       message.Content,
		ToolCalls:     make([]llm.ToolCall, len(message.ToolCalls)),
		ToolResponses: make([]llm.ToolResponse, len(message.ToolResponses)),
	}

	for iny, call := range message.ToolCalls {
		args, err := json.Marshal(call.Function.Arguments)
		if err != nil {
			return nil, err
		}

		result.ToolCalls[iny] = llm.ToolCall{
			CallID: call.CallID,
			Function: llm.Function{
				Name:      call.Function.Name,
				Arguments: string(args),
			},
		}
	}

	for iny, call := range message.ToolResponses {
		content, err := json.Marshal(call.Content)
		if err != nil {
			return nil, err
		}

		result.ToolResponses[iny] = llm.ToolResponse{
			CallID:  call.CallID,
			Content: string(content),
		}
	}

	return result, nil
}

// ToLLMMessages converts a slice of datastore messages to a slice of LLM messages.
func ToLLMMessages(message []*Message) ([]*llm.Message, error) {
	result := make([]*llm.Message, len(message))

	for inx, msg := range message {
		llmMsg, err := ToLLMMessage(msg)
		if err != nil {
			return nil, err
		}

		result[inx] = llmMsg
	}

	return result, nil
}

// ToDatastoreMessage converts a LLM message to a datastore message.
func ToDatastoreMessage(llmMessage *llm.Message) (*Message, error) {
	result := &Message{
		Role:          llmMessage.Role,
		Content:       llmMessage.Content,
		ToolCalls:     make([]ToolCall, len(llmMessage.ToolCalls)),
		ToolResponses: make([]ToolResponse, len(llmMessage.ToolResponses)),
	}

	for inx, call := range llmMessage.ToolCalls {
		args := make(map[string]interface{})
		err := json.Unmarshal([]byte(call.Function.Arguments), &args)
		if err != nil {
			return nil, err
		}

		result.ToolCalls[inx] = ToolCall{
			CallID: call.CallID,
			Function: Function{
				Name:      call.Function.Name,
				Arguments: args,
			},
		}
	}

	for inx, call := range llmMessage.ToolResponses {
		content := make(map[string]interface{})
		err := json.Unmarshal([]byte(call.Content), &content)
		if err != nil {
			return nil, err
		}

		result.ToolResponses[inx] = ToolResponse{
			CallID:  call.CallID,
			Content: content,
		}
	}

	return result, nil
}

// ToDatastoreMessages converts a slice of LLM messages to a slice of datastore messages.
func ToDatastoreMessages(llmMessage []*llm.Message) ([]*Message, error) {
	results := make([]*Message, len(llmMessage))

	for inx, msg := range llmMessage {
		datastoreMsg, err := ToDatastoreMessage(msg)
		if err != nil {
			return nil, err
		}

		results[inx] = datastoreMsg
	}

	return results, nil
}
