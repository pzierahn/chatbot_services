package chat

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"sort"
)

func getSources(messages []*llm.Message) []*pb.Source {
	if len(messages) < 3 {
		return nil
	}

	fragments := make(map[uuid.UUID][]*pb.Source_Fragment)

	for idx := len(messages) - 2; idx > 0; idx-- {
		message := messages[idx]
		if message.Role == llm.RoleUser && len(message.ToolResponses) == 0 {
			break
		}

		isSourceCall := make(map[string]bool)
		for _, toolCall := range messages[idx-1].ToolCalls {
			if toolCall.Name == "get_sources" {
				isSourceCall[toolCall.CallID] = true
			}
		}

		for _, toolResponse := range message.ToolResponses {
			if isSourceCall[toolResponse.CallID] {
				var source Sources

				err := json.Unmarshal([]byte(toolResponse.Content), &source)
				if err != nil {
					continue
				}

				for _, item := range source.Items {
					docId, err := uuid.Parse(item.DocumentId)
					if err != nil {
						continue
					}

					fragments[docId] = append(fragments[docId], &pb.Source_Fragment{
						Id:       item.Id,
						Content:  item.Text,
						Position: item.Position,
						Score:    item.Score,
					})
				}
			}
		}
	}

	sources := make([]*pb.Source, 0)
	for docId, parts := range fragments {
		sort.Slice(parts, func(i, j int) bool {
			return parts[i].Position < parts[j].Position
		})

		sources = append(sources, &pb.Source{
			DocumentId: docId.String(),
			Fragments:  parts,
		})
	}

	return sources
}

func messagesToProto(messages []*llm.Message) ([]*pb.Message, error) {
	protoMessages := make([]*pb.Message, 0)

	if len(messages)%2 != 0 {
		return nil, errors.New("invalid message count")
	}

	var idx int

	for {
		if idx >= len(messages) {
			break
		}

		user := messages[idx]
		protoMessage := &pb.Message{
			Prompt: user.Content,
		}

		assistant := messages[idx+1]

		var fragments = make(map[uuid.UUID][]*pb.Source_Fragment)

		for {
			assistant = messages[idx+1]
			if len(assistant.ToolCalls) == 0 {
				// Reached the end of tool calls
				break
			}

			isSourceCall := make(map[string]bool)
			for _, toolCall := range assistant.ToolCalls {
				if toolCall.Name == "get_sources" {
					isSourceCall[toolCall.CallID] = true
				}
			}

			user = messages[idx+2]
			for _, toolResponse := range user.ToolResponses {
				if !isSourceCall[toolResponse.CallID] {
					continue
				}

				var source Sources
				err := json.Unmarshal([]byte(toolResponse.Content), &source)
				if err != nil {
					return nil, err
				}

				for _, item := range source.Items {
					docId, err := uuid.Parse(item.DocumentId)
					if err != nil {
						continue
					}

					fragments[docId] = append(fragments[docId], &pb.Source_Fragment{
						Id:       item.Id,
						Content:  item.Text,
						Position: item.Position,
						Score:    item.Score,
					})
				}
			}

			idx += 2
		}

		protoMessage.Sources = make([]*pb.Source, 0)
		for docId, parts := range fragments {
			sort.Slice(parts, func(i, j int) bool {
				return parts[i].Position < parts[j].Position
			})

			protoMessage.Sources = append(protoMessage.Sources, &pb.Source{
				DocumentId: docId.String(),
				Fragments:  parts,
			})
		}

		protoMessage.Completion = assistant.Content
		protoMessages = append(protoMessages, protoMessage)

		idx += 2
	}

	return protoMessages, nil
}
