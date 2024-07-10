package qdrant

import (
	"crypto/tls"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/search"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

type Search struct {
	conn          *grpc.ClientConn
	apiKey        string
	namespace     string
	embedding     llm.Embedding
	fastEmbedding *search.ParallelEmbedding
	dimension     int
}

func (db *Search) Close() error {
	return db.conn.Close()
}

func New(engine llm.Embedding, namespace string) (*Search, error) {
	apiKey := os.Getenv("CHATBOT_QDRANT_KEY")
	target := os.Getenv("CHATBOT_QDRANT_URL")

	var opts []grpc.DialOption
	if os.Getenv("CHATBOT_QDRANT_INSECURE") == "true" {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, err
	}

	fastEmbedding := search.NewParallelEmbedding(engine, 10, 30)

	client := &Search{
		conn:          conn,
		apiKey:        apiKey,
		namespace:     namespace,
		embedding:     engine,
		dimension:     engine.GetEmbeddingDimension(),
		fastEmbedding: fastEmbedding,
	}

	err = client.Init()
	if err != nil {
		return nil, err
	}

	return client, nil
}
