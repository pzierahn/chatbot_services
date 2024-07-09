package notion

import (
	"context"
	"errors"
	"fmt"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
)

// InsertAPIKey inserts a new API key into the database.
func (client *Client) InsertAPIKey(ctx context.Context, key *pb.NotionApiKey) (*emptypb.Empty, error) {

	userId, err := client.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	// Insert the key into the database
	err = client.Database.InsertNotionAPIKey(ctx, &datastore.NotionAPIKey{
		UserId: userId,
		Key:    key.Key,
	})
	if err != nil {
		return nil, err
	}

	client.setCachedKey(userId, key.Key)

	return &emptypb.Empty{}, nil
}

// UpdateAPIKey updates an existing API key in the database.
func (client *Client) UpdateAPIKey(ctx context.Context, key *pb.NotionApiKey) (*emptypb.Empty, error) {

	userId, err := client.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	// Update the key in the database
	err = client.Database.UpdateNotionAPIKey(ctx, &datastore.NotionAPIKey{
		UserId: userId,
		Key:    key.Key,
	})
	if err != nil {
		return nil, err
	}

	client.setCachedKey(userId, key.Key)
	return &emptypb.Empty{}, nil
}

// DeleteAPIKey deletes an existing API key from the database.
func (client *Client) DeleteAPIKey(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {

	userId, err := client.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	// Delete the key from the database
	err = client.Database.DeleteNotionAPIKey(ctx, userId)
	if err != nil {
		return nil, err
	}

	client.deleteCachedKey(userId)
	return &emptypb.Empty{}, nil
}

// GetAPIKey retrieve the API key from the database.
func (client *Client) GetAPIKey(ctx context.Context, _ *emptypb.Empty) (*pb.NotionApiKey, error) {

	userId, err := client.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	key, ok := client.getCachedKey(userId)
	apiKey := &pb.NotionApiKey{}

	if ok {
		apiKey.Key = key
		return apiKey, nil
	}

	// Get the API key from the database
	apiKey.Key, err = client.Database.GetNotionAPIKey(ctx, userId)
	if errors.Is(err, mongo.ErrNoDocuments) {
		// No results found, return empty key
		apiKey.Key = ""
		return apiKey, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not get API key: %v", err)
	}

	client.setCachedKey(userId, apiKey.Key)
	return apiKey, nil
}
