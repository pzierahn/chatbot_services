package notion

import (
	"context"
	"github.com/jomei/notionapi"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

const fileColumn = "ID"

// ListDatabases retrieves all databases in the workspace.
func (client *Client) ListDatabases(ctx context.Context, _ *emptypb.Empty) (*pb.Databases, error) {
	resp, err := client.api.Search.Do(ctx, &notionapi.SearchRequest{
		Filter: notionapi.SearchFilter{
			Value:    "database",
			Property: "object",
		},
		PageSize: 999,
	})
	if err != nil {
		return nil, err
	}

	databases := &pb.Databases{}

	for _, result := range resp.Results {
		switch result.GetObject() {
		case "database":
			database := result.(*notionapi.Database)
			databases.Items = append(databases.Items, &pb.Databases_Item{
				Id:   database.ID.String(),
				Name: database.Title[0].PlainText,
			})
		}
	}

	return databases, nil
}

// CreateDatabase creates a new database in the workspace.
func (client *Client) CreateDatabase(ctx context.Context, title string) (*notionapi.Database, error) {
	resp, err := client.api.Database.Create(ctx, &notionapi.DatabaseCreateRequest{
		Parent: notionapi.Parent{
			Type:   "page_id",
			PageID: "9579f240b48c453b8af0bf129fb1881e",
		},
		Title: []notionapi.RichText{
			{
				Type: notionapi.ObjectTypeText,
				Text: &notionapi.Text{
					Content: title,
				},
			},
		},
		Properties: map[string]notionapi.PropertyConfig{
			fileColumn: notionapi.TitlePropertyConfig{
				Type: notionapi.PropertyConfigTypeTitle,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) ListDocumentIDs(ctx context.Context, databaseID string) (map[string]string, error) {
	dbEntries, err := client.api.Database.Query(
		ctx, notionapi.DatabaseID(databaseID),
		&notionapi.DatabaseQueryRequest{
			Sorts: []notionapi.SortObject{
				{
					Property:  fileColumn,
					Direction: "ascending",
				},
			},
			PageSize: 999,
		})
	if err != nil {
		return nil, err
	}

	pageIDs := make(map[string]string)
	for _, result := range dbEntries.Results {
		props := result.Properties

		rich, ok := props[fileColumn].(*notionapi.TitleProperty)
		if !ok {
			continue
		}

		if len(rich.Title) <= 0 {
			continue
		}

		pageIDs[rich.Title[0].PlainText] = result.ID.String()
	}

	return pageIDs, nil
}

func (client *Client) AddColumn(ctx context.Context, databaseID, title string) error {
	_, err := client.api.Database.Update(
		ctx,
		notionapi.DatabaseID(databaseID),
		&notionapi.DatabaseUpdateRequest{
			Properties: map[string]notionapi.PropertyConfig{
				title: notionapi.RichTextPropertyConfig{
					Type: notionapi.PropertyConfigTypeRichText,
				},
			},
		})

	return err
}

func (client *Client) UpdateRow(ctx context.Context, pageID, column, text string) error {
	_, err := client.api.Page.Update(
		ctx,
		notionapi.PageID(pageID),
		&notionapi.PageUpdateRequest{
			Properties: map[string]notionapi.Property{
				column: notionapi.RichTextProperty{
					Type: notionapi.PropertyTypeRichText,
					RichText: []notionapi.RichText{
						{
							Text: &notionapi.Text{
								Content: text,
							},
						},
					},
				},
			},
		})

	return err
}
