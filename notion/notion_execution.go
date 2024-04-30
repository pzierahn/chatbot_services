package notion

import (
	"github.com/pzierahn/chatbot_services/llm/bedrock"
	pb "github.com/pzierahn/chatbot_services/proto"
	"sync"
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
		return err
	}

	names, err := client.documents.MapDocumentNames(ctx, &pb.CollectionID{
		Id: prompt.CollectionID,
	})

	err = client.AddColumn(ctx, prompt.DatabaseID, prompt.Prompt)
	if err != nil {
		return err
	}

	jobs := make([]func() (string, error), 0)

	for documentName, pageID := range pageIDs {
		documentID, ok := names.Items[documentName]
		if !ok {
			continue
		}

		jobs = append(jobs, func() (string, error) {
			resp, err := client.chat.Completion(ctx, &pb.CompletionRequest{
				DocumentId:   documentID,
				Prompt:       prompt.Prompt,
				ModelOptions: model,
			})
			if err != nil {
				return documentName, err
			}

			err = client.UpdateRow(ctx, pageID, prompt.Prompt, resp.Completion)
			if err != nil {
				return documentName, err
			}

			return documentName, nil
		})
	}

	var wg sync.WaitGroup
	wg.Add(len(jobs))

	agents := 10
	mux := make(chan struct{}, agents)
	for range agents {
		mux <- struct{}{}
	}

	for _, job := range jobs {
		go func(job func() (string, error)) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			case _, ok := <-mux:
				if !ok {
					return
				}
			}

			defer func() { mux <- struct{}{} }()

			documentName, err := job()
			if err != nil {
				return
			}

			_ = stream.Send(&pb.ExecutionResult{
				Document: documentName,
			})
		}(job)
	}

	wg.Wait()

	return nil
}
