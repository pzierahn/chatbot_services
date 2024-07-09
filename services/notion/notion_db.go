package notion

import (
	"context"
	"github.com/jomei/notionapi"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

const fileColumn = "ID"

// ListDatabases retrieves all databases in the workspace.
func (client *Client) ListDatabases(ctx context.Context, _ *emptypb.Empty) (*pb.Databases, error) {
	api, err := client.getAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := api.Search.Do(ctx, &notionapi.SearchRequest{
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

// filenamesPageIds maps filenames to page IDs in a database.
func (client *Client) filenamesPageIds(ctx context.Context, databaseId string) (map[string]string, error) {
	api, err := client.getAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	dbEntries, err := api.Database.Query(
		ctx, notionapi.DatabaseID(databaseId),
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

		filename := rich.Title[0].PlainText
		filename = strings.TrimSuffix(filename, ".pdf")

		pageIDs[filename] = result.ID.String()
	}

	return pageIDs, nil
}

// addNewColumn adds a new column with the given title to a database.
func (client *Client) addNewColumn(ctx context.Context, databaseID, title string) error {
	api, err := client.getAPIClient(ctx)
	if err != nil {
		return err
	}

	_, err = api.Database.Update(
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
	api, err := client.getAPIClient(ctx)
	if err != nil {
		return err
	}

	_, err = api.Page.Update(
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
