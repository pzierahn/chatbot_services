package table

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/llm/bedrock"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/utils"
	"log"
	"time"
)

type Column struct {
	ID               string `json:"id"`
	GenerationPrompt string `json:"prompt"`
	Name             string `json:"name"`
}

type Cell struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type Row struct {
	ID    string `json:"id"`
	Cells []Cell `json:"cells"`
}

type Table struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Columns   []Column  `json:"columns"`
	Rows      []Row     `json:"rows"`
}

func docToText(doc *pb.Document) string {
	var text string

	for _, chunk := range doc.Chunks {
		text += chunk.Text
	}

	return text
}

func (service *Service) CreateTable(ctx context.Context, req *pb.NewTable) (*pb.TableId, error) {
	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	table := Table{
		ID:        req.Id,
		Name:      req.Name,
		CreatedAt: time.Now(),
		Columns:   []Column{},
		Rows:      []Row{},
	}

	for _, column := range req.Columns {
		log.Printf("Generating column name from prompt: '%s'", column)

		resp, err := service.agent.GenerateCompletion(ctx, &llm.GenerateRequest{
			SystemPrompt: "Don't write any introductions, preambles, or other text that isn't part of the table. Just write the column name.",
			Messages: []*llm.Message{{
				Type: llm.MessageTypeUser,
				Text: fmt.Sprintf("Create a table column name from this prompt in snake case: '%s'. Do it without any additional words or explanations.", column),
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

		table.Columns = append(table.Columns, Column{
			ID:               uuid.NewString(),
			GenerationPrompt: column,
			Name:             resp.Text,
		})
	}

	ch := make(chan Row, len(req.DocumentIds))

	for _, docID := range req.DocumentIds {
		go func(docID string) {
			log.Printf("Generating row from document: %s", docID)

			doc, err := service.document.Get(ctx, &pb.DocumentID{Id: docID})
			if err != nil {
				log.Printf("Error getting document: %s", err)
				return
			}

			row := Row{
				ID:    uuid.NewString(),
				Cells: []Cell{},
			}

			for _, column := range table.Columns {
				text := docToText(doc)

				resp, err := service.agent.GenerateCompletion(ctx, &llm.GenerateRequest{
					SystemPrompt: "Return only the requested information, without any additional words or explanation. " +
						"The answer should be as short as possible. " +
						"Don't repeat the question or prompt.",
					Messages: []*llm.Message{{
						Type: llm.MessageTypeUser,
						Text: text + "\n\n" + column.GenerationPrompt,
					}},
					Model:       bedrock.ClaudeSonnet,
					MaxTokens:   100,
					TopP:        0,
					Temperature: 0,
					UserId:      userId,
				})
				if err != nil {
					log.Printf("Error generating completion: %s", err)
					return
				}

				cell := Cell{
					ID:    column.ID,
					Value: resp.Text,
				}

				row.Cells = append(row.Cells, cell)
			}

			ch <- row
		}(docID)
	}

	for row := range ch {
		table.Rows = append(table.Rows, row)

		if len(table.Rows) == len(req.DocumentIds) {
			close(ch)
		}
	}

	log.Printf("Creating table: %s", utils.Prettify(table))
	utils.WriteJson(table, "table.json")

	var id string
	err = service.db.QueryRow(ctx,
		`INSERT INTO user_tables (user_id, name, data)
			VALUES ($1, $2, $3)
            RETURNING id`,
		userId, table.Name, table).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &pb.TableId{Id: id}, nil
}
