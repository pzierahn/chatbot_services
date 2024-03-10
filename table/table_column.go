package table

import (
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/llm/bedrock"
	pb "github.com/pzierahn/chatbot_services/proto"
)

// AddColumnToTable adds a new column to a table.
func (service *Service) AddColumnToTable(ctx context.Context, req *pb.NewColumn) (*pb.Column, error) {
	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := service.agent.GenerateCompletion(ctx, &llm.GenerateRequest{
		Messages: []*llm.Message{{
			Type: llm.MessageTypeUser,
			Text: fmt.Sprintf("Create a table column name from this prompt in snake case: '%s'. "+
				"Do it without any additional words or explanations.", req.GenerationPrompt),
		}},
		Model:       bedrock.ClaudeSonnet,
		MaxTokens:   10,
		TopP:        0,
		Temperature: 0,
		UserId:      userId,
	})
	if err != nil {
		return nil, err
	}

	name := resp.Text

	var id string
	err = service.db.QueryRow(ctx,
		`INSERT INTO user_table_columns(user_id, table_id, generation_prompt, name)
			VALUES ($1, $2, $3, $4)
            RETURNING id`,
		userId, req.TableId, req.GenerationPrompt, name).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &pb.Column{
		Id:               id,
		Name:             name,
		GenerationPrompt: req.GenerationPrompt,
	}, nil
}
