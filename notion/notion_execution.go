package notion

import (
	"github.com/pzierahn/chatbot_services/llm/bedrock"
	pb "github.com/pzierahn/chatbot_services/proto"
	"log"
)

var model = &pb.ModelOptions{
	Model:       bedrock.ClaudeHaiku,
	Temperature: 1.0,
	MaxTokens:   256,
	TopP:        1.0,
}

func (client *Client) ExecutePrompt(prompt *pb.NotionPrompt, stream pb.Notion_ExecutePromptServer) error {
	ctx := stream.Context()

	pageIDs, err := client.ListDocumentIDs(ctx, prompt.DatabaseID)
	if err != nil {
		log.Println(err)
		return err
	}

	names, err := client.documents.MapDocumentNames(ctx, &pb.CollectionID{
		Id: prompt.CollectionID,
	})
	log.Printf("Names: %s", names)

	err = client.AddColumn(ctx, prompt.DatabaseID, prompt.Prompt)
	if err != nil {
		log.Println(err)
		return err
	}

	for documentName, pageID := range pageIDs {
		documentID := names.Items[documentName]
		log.Printf("Document: %s (%s)", documentName, documentID)

		resp, err := client.chat.Completion(ctx, &pb.CompletionRequest{
			DocumentId:   documentID,
			Prompt:       prompt.Prompt,
			ModelOptions: model,
		})
		if err != nil {
			log.Println(err)
			return err
		}

		err = client.UpdateRow(ctx, pageID, prompt.Prompt, resp.Completion)
		if err != nil {
			log.Println(err)
			return err
		}

		_ = stream.Send(&pb.ExecutionResult{
			Document: documentName,
		})
	}

	return nil
}
