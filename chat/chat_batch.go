package chat

import (
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"log"
)

type batchJob struct {
	docIdx    uint32
	promptIdx uint32
	prompt    string
	text      string
	title     string
	retry     int
}

type batchResults struct {
	job        batchJob
	completion string
	err        error
	retry      int
}

func getDocumentText(doc *pb.Document) string {
	var text string

	for _, chunk := range doc.Chunks {
		text += chunk.Text
	}

	return text
}

func getDocumentTitle(doc *pb.Document) string {
	switch doc.Metadata.Data.(type) {
	case *pb.DocumentMetadata_Web:
		return doc.Metadata.GetWeb().Title
	case *pb.DocumentMetadata_File:
		return doc.Metadata.GetFile().Filename
	}

	return ""
}

func (service *Service) BatchChat(ctx context.Context, req *pb.BatchRequest) (*pb.BatchResponse, error) {
	userId, err := service.Verify(ctx)
	if err != nil {
		return nil, err
	}

	model, err := service.getModel(req.ModelOptions.Model)
	if err != nil {
		return nil, err
	}

	response := &pb.BatchResponse{
		DocumentIds: req.DocumentIds,
		Prompts:     req.Prompts,
		PromptTitle: make([]string, len(req.Prompts)),
	}

	for promptIdx, prompt := range req.Prompts {
		resp, err := model.GenerateCompletion(ctx, &llm.GenerateRequest{
			Messages: []*llm.Message{{
				Type: llm.MessageTypeUser,
				Text: fmt.Sprintf("Create a table column name from this prompt in snake case: '%s'. "+
					"Do it without any additional words or explanations.", prompt),
			}},
			Model:       req.ModelOptions.Model,
			MaxTokens:   10,
			TopP:        0,
			Temperature: 0,
			UserId:      userId,
		})
		if err != nil {
			return nil, err
		}

		response.PromptTitle[promptIdx] = resp.Text
	}

	scheduledJobs := 0
	jobQueue := make(chan batchJob, len(req.DocumentIds))
	defer close(jobQueue)

	for docIdx, docID := range req.DocumentIds {
		doc, err := service.docs.Get(ctx, &pb.DocumentID{Id: docID})
		if err != nil {
			return nil, err
		}

		text := getDocumentText(doc)
		title := getDocumentTitle(doc)

		for promptIdx, prompt := range req.Prompts {
			jobQueue <- batchJob{
				docIdx:    uint32(docIdx),
				promptIdx: uint32(promptIdx),
				prompt:    prompt,
				text:      text,
				title:     title,
			}
			scheduledJobs++
		}
	}

	resultsQueue := make(chan batchResults, 4)
	defer close(resultsQueue)

	for worker := 0; worker < 10; worker++ {
		go func() {
			for job := range jobQueue {
				completion, cerr := model.GenerateCompletion(ctx, &llm.GenerateRequest{
					Model:       req.ModelOptions.Model,
					MaxTokens:   int(req.ModelOptions.MaxTokens),
					Temperature: req.ModelOptions.Temperature,
					TopP:        req.ModelOptions.TopP,
					UserId:      userId,
					Messages: []*llm.Message{{
						Text: job.text + "\n\n\n" + job.prompt,
					}},
				})
				if cerr != nil {
					resultsQueue <- batchResults{
						job:   job,
						err:   cerr,
						retry: job.retry + 1,
					}

					log.Println(cerr)
				} else {
					resultsQueue <- batchResults{
						job:        job,
						completion: completion.Text,
					}
				}
			}
		}()
	}

	for result := range resultsQueue {
		if result.err != nil {
			if result.retry < 2 {
				jobQueue <- result.job
			} else {
				return nil, result.err
			}

			continue
		}

		response.Items = append(response.Items, &pb.BatchResponse_Completion{
			DocumentId:    result.job.docIdx,
			DocumentTitle: result.job.title,
			Prompt:        result.job.promptIdx,
			Completion:    result.completion,
		})

		if scheduledJobs == len(response.Items) {
			break
		}
	}

	return response, nil
}
