package test

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

type Tester struct {
	chat pb.ChatServiceClient
}

func (test Tester) runTest(name string, testFunc func(ctx context.Context) error) {
	uid := uuid.NewString()

	ctx, cnl := context.WithTimeout(context.Background(), time.Second*10)
	defer cnl()

	ctx = metadata.AppendToOutgoingContext(ctx, "User-Id", uid)

	start := time.Now()
	err := testFunc(ctx)
	execTime := time.Since(start)

	if err != nil {
		log.Printf("[%v] FAIL: %v", name, err)
	} else {
		log.Printf("[%v] PASS: %v", name, execTime)
	}
}

func (test Tester) expectError(name string, testFunc func(ctx context.Context) error) {
	uid := uuid.NewString()
	ctx, cnl := context.WithTimeout(context.Background(), time.Second*10)
	defer cnl()

	ctx = metadata.AppendToOutgoingContext(ctx, "User-Id", uid)

	start := time.Now()
	err := testFunc(ctx)
	execTime := time.Since(start)

	if err == nil {
		log.Printf("[%v] FAIL: expected error", name)
	} else {
		log.Printf("[%v] PASS: %v", name, execTime)
	}
}

func NewTester(chat pb.ChatServiceClient) *Tester {
	return &Tester{
		chat: chat,
	}
}
