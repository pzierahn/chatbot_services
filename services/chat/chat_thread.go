package chat

import (
	"context"
	"errors"
	"github.com/google/uuid"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ListThreadIDs returns a list of thread IDs for a given collection.
func (service *Service) ListThreadIDs(ctx context.Context, collection *pb.CollectionId) (*pb.ThreadIDs, error) {
	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	collectionId, err := uuid.Parse(collection.Id)
	if err != nil {
		return nil, err
	}

	threads, err := service.Database.GetThreadIDs(ctx, userId, collectionId)
	if err != nil {
		return nil, err
	}

	results := &pb.ThreadIDs{}
	for _, thread := range threads {
		results.Ids = append(results.Ids, thread.String())
	}

	return results, nil
}

// GetThread returns a thread by ID.
func (service *Service) GetThread(ctx context.Context, req *pb.ThreadID) (*pb.Thread, error) {
	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	threadId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	thread, err := service.Database.GetThread(ctx, userId, threadId)
	if err != nil {
		return nil, err
	}

	messages, err := messagesToProto(thread.Messages)
	if err != nil {
		return nil, err
	}

	results := &pb.Thread{
		Id:       req.Id,
		Messages: messages,
	}

	return results, nil
}

// DeleteThread deletes a thread by ID.
func (service *Service) DeleteThread(ctx context.Context, req *pb.ThreadID) (*emptypb.Empty, error) {
	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	threadId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	err = service.Database.DeleteThread(ctx, userId, threadId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// DeleteMessageFromThread deletes a message from a thread by index.
func (service *Service) DeleteMessageFromThread(ctx context.Context, req *pb.MessageIndex) (*emptypb.Empty, error) {
	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	threadId, err := uuid.Parse(req.ThreadId)
	if err != nil {
		return nil, err
	}

	// Get the thread for the message
	thread, err := service.Database.GetThread(ctx, userId, threadId)
	if err != nil {
		return nil, err
	}

	// Check if the message index is valid
	if req.Index >= uint32(len(thread.Messages)) {
		return nil, errors.New("invalid message index")
	}

	// Delete the message at position req.Index
	thread.Messages = append(thread.Messages[:req.Index], thread.Messages[req.Index+1:]...)

	// Store the thread back to the database
	err = service.Database.StoreThread(ctx, thread)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
