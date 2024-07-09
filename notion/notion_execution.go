package notion

import (
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/proto"
	"sync"
)

func (client *Client) ExecutePrompt(prompt *pb.NotionPrompt, stream pb.Notion_ExecutePromptServer) error {
	ctx := stream.Context()

	userId, err := client.Auth.VerifyFunding(ctx)
	if err != nil {
		return err
	}

	collectionId, err := uuid.Parse(prompt.CollectionId)
	if err != nil {
		return err
	}

	pageIDs, err := client.filenamesPageIds(ctx, prompt.DatabaseId)
	if err != nil {
		return err
	}

	err = client.addNewColumn(ctx, prompt.DatabaseId, prompt.Prompt)
	if err != nil {
		return err
	}

	jobs := make([]func() (string, error), 0)

	for documentName, pageID := range pageIDs {
		jobs = append(jobs, func() (string, error) {
			documentId, err := client.Database.FindDocumentId(ctx, datastore.FindRequest{
				UserId:       userId,
				CollectionId: collectionId,
				Name:         documentName,
			})
			if err != nil {
				return documentName, err
			}

			resp, err := client.Chat.Completion(ctx, &pb.CompletionRequest{
				DocumentId:   documentId.String(),
				Prompt:       prompt.Prompt,
				ModelOptions: prompt.ModelOptions,
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
