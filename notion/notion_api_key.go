package notion

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

var mux sync.RWMutex
var apiKeyCache = make(map[string]string)

func (client *Client) InsertAPIKey(ctx context.Context, key *pb.NotionApiKey) (*emptypb.Empty, error) {

	userId, err := client.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	// Insert the key into the database
	err = client.Database.InsertNotionAPIKey(ctx, &datastore.NotionAPIKey{
		Id:     uuid.New(),
		UserId: userId,
		Key:    key.Key,
	})
	if err != nil {
		return nil, err
	}

	mux.Lock()
	apiKeyCache[userId] = key.Key
	mux.Unlock()

	return &emptypb.Empty{}, nil
}

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

	mux.Lock()
	apiKeyCache[userId] = key.Key
	mux.Unlock()

	return &emptypb.Empty{}, nil
}

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

	mux.Lock()
	delete(apiKeyCache, userId)
	mux.Unlock()

	return &emptypb.Empty{}, nil
}

func (client *Client) GetAPIKey(ctx context.Context, _ *emptypb.Empty) (*pb.NotionApiKey, error) {

	userId, err := client.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	mux.RLock()
	key, ok := apiKeyCache[userId]
	mux.RUnlock()

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

	mux.Lock()
	apiKeyCache[userId] = apiKey.Key
	mux.Unlock()

	return apiKey, nil
}
